#!/bin/bash

# Version management for multi-module releases
# Handles version detection, validation, and selection

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"
# shellcheck source=./lib/config.sh
source "$SCRIPT_DIR/lib/config.sh"

# Retracted versions that should never be suggested
readonly RETRACTED_VERSIONS=("v1.0.0" "v1.0.1")

# Get current version for a specific module
get_module_version() {
    local tag_prefix="$1"

    local pattern
    if [[ "$tag_prefix" == "." ]]; then
        # Main module: match v*.*.* but not paths
        pattern="v[0-9]*.[0-9]*.[0-9]*"
        git tag -l "$pattern" | grep -v "/" | sort -V | tail -n1
    else
        # Submodule: match prefix/v*.*.*
        pattern="${tag_prefix}/v[0-9]*.[0-9]*.[0-9]*"
        git tag -l "$pattern" | sort -V | tail -n1 | sed "s|^${tag_prefix}/||"
    fi
}

# Get the latest version across all modules
get_latest_version() {
    local latest=""

    for entry in "${RELEASABLE_MODULES[@]}"; do
        local module_path tag_prefix
        read -r module_path tag_prefix <<< "$(parse_module_entry "$entry")"

        local version
        version=$(get_module_version "$tag_prefix")

        if [[ -n "$version" ]]; then
            if [[ -z "$latest" ]]; then
                latest="$version"
            else
                local cmp
                cmp=$(version_compare "$version" "$latest")
                if [[ "$cmp" == "gt" ]]; then
                    latest="$version"
                fi
            fi
        fi
    done

    echo "$latest"
}

# Check if a version is retracted
is_version_retracted() {
    local version="$1"

    for retracted in "${RETRACTED_VERSIONS[@]}"; do
        if [[ "$version" == "$retracted" ]]; then
            return 0
        fi
    done

    return 1
}

# Check if version exists for any module
version_exists() {
    local version="$1"

    for entry in "${RELEASABLE_MODULES[@]}"; do
        local module_path tag_prefix
        read -r module_path tag_prefix <<< "$(parse_module_entry "$entry")"

        local tag
        tag=$(get_module_tag "$tag_prefix" "$version")

        if git tag -l "$tag" | grep -q "^${tag}$"; then
            return 0
        fi
    done

    return 1
}

# Propose next versions (patch, minor, major)
propose_versions() {
    local current="$1"

    local patch minor major
    patch=$(increment_version "$current" "patch")
    minor=$(increment_version "$current" "minor")
    major=$(increment_version "$current" "major")

    # Check for retracted versions and adjust
    while is_version_retracted "$patch"; do
        patch=$(increment_version "$patch" "patch")
    done

    while is_version_retracted "$minor"; do
        minor=$(increment_version "$minor" "minor")
    done

    while is_version_retracted "$major"; do
        major=$(increment_version "$major" "major")
    done

    # Special case: if current is v0.x.x and major would be v1.0.0 (retracted)
    # suggest v2.0.0 instead
    if [[ "$major" == "v1.0.0" ]]; then
        major="v2.0.0"
    fi

    echo "$patch|$minor|$major"
}

# Show current versions for all modules
show_versions() {
    log_section "Current Module Versions"

    local max_version=""

    for entry in "${RELEASABLE_MODULES[@]}"; do
        local module_path tag_prefix
        read -r module_path tag_prefix <<< "$(parse_module_entry "$entry")"

        local display_path="$module_path"
        [[ "$display_path" == "." ]] && display_path="(main)"

        local version
        version=$(get_module_version "$tag_prefix")

        if [[ -z "$version" ]]; then
            echo "  $display_path: ${COLOR_GRAY}(no version)${COLOR_RESET}"
        else
            echo "  $display_path: $version"

            # Track highest version
            if [[ -z "$max_version" ]] || [[ $(version_compare "$version" "$max_version") == "gt" ]]; then
                max_version="$version"
            fi
        fi
    done

    echo ""
    if [[ -n "$max_version" ]]; then
        log_info "Latest version across all modules: ${COLOR_GREEN}$max_version${COLOR_RESET}"
    fi
}

# Interactive version selection
interactive_version_select() {
    local current_version
    current_version=$(get_latest_version)

    if [[ -z "$current_version" ]]; then
        current_version="v0.0.0"
        log_warning "No existing versions found, starting from $current_version"
    fi

    log_section "Version Selection"

    echo "Current latest version: ${COLOR_GREEN}$current_version${COLOR_RESET}"
    echo ""

    # Get proposed versions
    local proposals
    proposals=$(propose_versions "$current_version")
    IFS='|' read -r patch_version minor_version major_version <<< "$proposals"

    echo "Select version for release:"
    echo "  ${COLOR_GREEN}1)${COLOR_RESET} $patch_version (patch) - Bug fixes"
    echo "  ${COLOR_GREEN}2)${COLOR_RESET} $minor_version (minor) - New features [recommended]"
    echo "  ${COLOR_GREEN}3)${COLOR_RESET} $major_version (major) - Breaking changes"
    echo "  ${COLOR_GREEN}4)${COLOR_RESET} Custom version"
    echo ""

    local choice
    choice=$(prompt "Choice" "2")

    local selected_version
    case "$choice" in
        1)
            selected_version="$patch_version"
            ;;
        2)
            selected_version="$minor_version"
            ;;
        3)
            selected_version="$major_version"
            ;;
        4)
            while true; do
                selected_version=$(prompt "Enter version (format: vX.Y.Z)")

                if ! is_semantic_version "$selected_version"; then
                    log_error "Invalid version format. Must be vX.Y.Z (e.g., v0.19.0)"
                    continue
                fi

                if is_version_retracted "$selected_version"; then
                    log_error "Version $selected_version is retracted and cannot be used"
                    log_info "Retracted versions: ${RETRACTED_VERSIONS[*]}"
                    continue
                fi

                break
            done
            ;;
        *)
            die "Invalid choice: $choice"
            ;;
    esac

    # Validate selected version
    if is_version_retracted "$selected_version"; then
        die "Version $selected_version is retracted (see go.mod). Please choose a different version."
    fi

    if version_exists "$selected_version"; then
        die "Version $selected_version already exists. Please choose a different version."
    fi

    echo ""
    log_success "Selected version: ${COLOR_GREEN}$selected_version${COLOR_RESET}"
    echo ""

    echo "$selected_version"
}

# Main script
main() {
    local mode="${1:-}"

    case "$mode" in
        --show)
            show_versions
            ;;
        --interactive)
            interactive_version_select
            ;;
        --latest)
            get_latest_version
            ;;
        --validate)
            shift
            local version="$1"
            if is_semantic_version "$version"; then
                if is_version_retracted "$version"; then
                    log_error "Version $version is retracted"
                    exit 1
                fi
                if version_exists "$version"; then
                    log_error "Version $version already exists"
                    exit 1
                fi
                log_success "Version $version is valid"
                exit 0
            else
                log_error "Invalid version format: $version"
                exit 1
            fi
            ;;
        --help|-h)
            cat <<EOF
Usage: version.sh [OPTIONS]

Version management for Fuego multi-module releases.

OPTIONS:
    --show              Show current versions for all modules
    --interactive       Interactive version selection (default)
    --latest            Print latest version across all modules
    --validate VERSION  Validate a version format and availability
    --help, -h          Show this help message

EXAMPLES:
    # Show current versions
    ./version.sh --show

    # Interactive selection
    ./version.sh --interactive

    # Check latest version
    ./version.sh --latest

    # Validate version
    ./version.sh --validate v0.19.0
EOF
            ;;
        "")
            # Default: interactive
            interactive_version_select
            ;;
        *)
            log_error "Unknown option: $mode"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
}

# Run main if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi

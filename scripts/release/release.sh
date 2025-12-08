#!/bin/bash

# Main orchestrator for Fuego multi-module releases
# Coordinates validation, version selection, and tagging

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"
# shellcheck source=./lib/config.sh
source "$SCRIPT_DIR/lib/config.sh"

# Parse command line arguments
DRY_RUN=0
VERSION=""
SKIP_VALIDATION=0

parse_args() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --dry-run)
                DRY_RUN=1
                shift
                ;;
            --version)
                VERSION="$2"
                shift 2
                ;;
            --skip-validation)
                SKIP_VALIDATION=1
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done
}

show_help() {
    cat <<EOF
Usage: release.sh [OPTIONS]

Main orchestrator for Fuego multi-module releases.

OPTIONS:
    --dry-run              Preview mode - no tags created or pushed
    --version VERSION      Use specific version (skip interactive selection)
    --skip-validation      Skip validation checks (not recommended)
    --help, -h             Show this help message

ENVIRONMENT VARIABLES:
    FUEGO_RELEASE_AUTO_CONFIRM=1    Skip all confirmation prompts
    FUEGO_RELEASE_SKIP_LINT=1       Skip linting
    FUEGO_RELEASE_SKIP_TESTS=1      Skip tests
    FUEGO_RELEASE_SKIP_BUILD=1      Skip builds
    FUEGO_RELEASE_ALLOW_NON_MAIN=1  Allow release from non-main branch
    DEBUG=1                          Enable debug logging

EXAMPLES:
    # Interactive release
    ./release.sh

    # Dry run (preview only)
    ./release.sh --dry-run

    # Release specific version
    ./release.sh --version v0.19.0

    # Quick release (skip validation - risky!)
    ./release.sh --skip-validation

MAKEFILE SHORTCUTS:
    make release              # Interactive release
    make release-dry-run      # Dry run mode
    make release-versions     # Show current versions
    make release-validate     # Run validation only
    make release-rollback     # Emergency rollback
EOF
}

# Show current state
show_current_state() {
    log_section "Current State"

    local current_branch
    current_branch=$(get_current_branch)

    echo "  Repository: github.com/go-fuego/fuego"
    echo "  Branch: $current_branch"
    echo "  Modules: $(count_modules) releasable modules"

    local latest_version
    latest_version=$("$SCRIPT_DIR/version.sh" --latest)

    if [[ -n "$latest_version" ]]; then
        echo "  Latest release: $latest_version"
    else
        echo "  Latest release: (none)"
    fi

    echo ""
}

# Show release summary
show_release_summary() {
    local version="$1"

    log_section "Release Summary"

    echo "Version: ${COLOR_GREEN}$version${COLOR_RESET}"
    echo "Modules: $(count_modules) modules (all at same version)"
    echo ""
    echo "Tags to create:"

    for entry in "${RELEASABLE_MODULES[@]}"; do
        local module_path tag_prefix
        read -r module_path tag_prefix <<< "$(parse_module_entry "$entry")"

        local tag
        tag=$(get_module_tag "$tag_prefix" "$version")

        echo "  • $tag"
    done

    echo ""
}

# Main release workflow
main() {
    parse_args "$@"

    if [[ "$DRY_RUN" == "1" ]]; then
        print_banner "Fuego Multi-Module Release (DRY RUN)"
        log_warning "DRY RUN MODE - No actual changes will be made"
        echo ""
    else
        print_banner "Fuego Multi-Module Release"
    fi

    # Show current state
    show_current_state

    # Step 1: Run pre-release validation
    if [[ "$SKIP_VALIDATION" == "1" ]]; then
        log_warning "Skipping validation (--skip-validation)"
    else
        if [[ "$DRY_RUN" == "1" ]]; then
            log_info "Skipping validation in dry-run mode"
        else
            log_info "Running pre-release checks..."
            echo ""

            if ! "$SCRIPT_DIR/validate.sh" full; then
                die "Pre-release validation failed. Fix issues and try again."
            fi

            echo ""
        fi
    fi

    # Step 2: Version selection
    if [[ -z "$VERSION" ]]; then
        if [[ "$DRY_RUN" == "1" ]]; then
            # For dry-run, use a test version
            local latest
            latest=$("$SCRIPT_DIR/version.sh" --latest 2>/dev/null || echo "v0.18.8")
            # Extract major.minor.patch and increment minor
            local major minor patch
            latest="${latest#v}"
            IFS='.' read -r major minor patch <<< "$latest"
            minor=$((minor + 1))
            VERSION="v${major}.${minor}.0"
            log_info "Using test version for dry-run: $VERSION"
            echo ""
        else
            VERSION=$("$SCRIPT_DIR/version.sh" --interactive)
        fi
    else
        log_info "Using specified version: $VERSION"

        # Validate the provided version
        if [[ "$DRY_RUN" != "1" ]]; then
            if ! "$SCRIPT_DIR/version.sh" --validate "$VERSION" >/dev/null 2>&1; then
                die "Invalid version: $VERSION"
            fi
        fi

        echo ""
    fi

    # Step 3: Show summary and confirm
    show_release_summary "$VERSION"

    if [[ "$DRY_RUN" == "1" ]]; then
        log_info "[DRY RUN] Would create and push these tags"
        echo ""

        # Preview tags
        "$SCRIPT_DIR/tag.sh" --preview "$VERSION"
        echo ""

        log_success "[DRY RUN] Complete - no actual changes made"
        echo ""
        log_info "To perform actual release, run without --dry-run"
        exit 0
    fi

    # Confirm before proceeding
    if ! confirm "Proceed with release?" "n"; then
        log_info "Release cancelled"
        exit 0
    fi

    echo ""

    # Step 4: Create tags
    log_info "Creating tags..."
    echo ""

    if ! "$SCRIPT_DIR/tag.sh" --create "$VERSION"; then
        die "Tag creation failed"
    fi

    echo ""

    # Step 5: Push tags
    if ! confirm "Push tags to remote (github.com/go-fuego/fuego)?" "n"; then
        log_warning "Tags created locally but not pushed"
        log_info "To push later: $SCRIPT_DIR/tag.sh --push"
        log_info "To rollback: $SCRIPT_DIR/tag.sh --rollback"
        exit 0
    fi

    echo ""

    if ! "$SCRIPT_DIR/tag.sh" --push; then
        log_error "Failed to push tags"
        log_info "Tags remain local. You can:"
        log_info "  1. Fix the issue and retry: $SCRIPT_DIR/tag.sh --push"
        log_info "  2. Rollback: $SCRIPT_DIR/tag.sh --rollback"
        exit 1
    fi

    # Step 6: Success!
    echo ""
    log_section "Release Complete!"
    log_success "Fuego $VERSION has been released!"
    echo ""

    echo "Next steps:"
    echo "  1. Verify tags on GitHub:"
    echo "     https://github.com/go-fuego/fuego/tags"
    echo ""
    echo "  2. (Optional) Create GitHub releases with:"
    echo "     gh release create $VERSION --generate-notes"
    echo ""
    echo "  3. (Optional) Update submodule dependencies in go.mod files"
    echo "     to reference $VERSION"
    echo ""
    echo "  4. Announce the release!"
    echo ""

    log_success "All done!"
}

# Run main
main "$@"

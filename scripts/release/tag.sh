#!/bin/bash

# Tag management for multi-module releases
# Creates, pushes, and manages git tags atomically

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"
# shellcheck source=./lib/config.sh
source "$SCRIPT_DIR/lib/config.sh"

# Track created tags for rollback
declare -a CREATED_TAGS=()

# Store tags in a temp file for persistence across script calls
TAGS_FILE="/tmp/fuego-release-tags-$$.txt"

# Save created tags to file
save_created_tags() {
    printf '%s\n' "${CREATED_TAGS[@]}" > "$TAGS_FILE"
}

# Load created tags from file
load_created_tags() {
    if [[ -f "$TAGS_FILE" ]]; then
        mapfile -t CREATED_TAGS < "$TAGS_FILE"
    fi
}

# Create a single tag
create_tag() {
    local tag="$1"
    local message="${2:-Release $tag}"
    local dry_run="${3:-0}"

    if [[ "$dry_run" == "1" ]]; then
        echo "  [DRY RUN] Would create: $tag"
        return 0
    fi

    if git tag -a "$tag" -m "$message" 2>/dev/null; then
        log_success "Created $tag"
        CREATED_TAGS+=("$tag")
        return 0
    else
        log_error "Failed to create tag: $tag"
        return 1
    fi
}

# Preview tags that will be created
preview_tags() {
    local version="$1"

    log_section "Tags to Create"

    echo "Version: ${COLOR_GREEN}$version${COLOR_RESET}"
    echo "Number of tags: $(count_modules)"
    echo ""

    local index=1
    for entry in "${RELEASABLE_MODULES[@]}"; do
        local module_path tag_prefix
        read -r module_path tag_prefix <<< "$(parse_module_entry "$entry")"

        local tag
        tag=$(get_module_tag "$tag_prefix" "$version")

        local display_path="$module_path"
        [[ "$display_path" == "." ]] && display_path="(main)"

        printf "  %d. %-40s %s\n" "$index" "$tag" "$display_path"
        index=$((index + 1))
    done

    echo ""
}

# Create all tags
create_all_tags() {
    local version="$1"
    local dry_run="${2:-0}"

    if [[ "$dry_run" == "1" ]]; then
        log_section "Creating Tags (DRY RUN)"
    else
        log_section "Creating Tags"
    fi

    local success=true

    for entry in "${RELEASABLE_MODULES[@]}"; do
        local module_path tag_prefix
        read -r module_path tag_prefix <<< "$(parse_module_entry "$entry")"

        local tag
        tag=$(get_module_tag "$tag_prefix" "$version")

        if ! create_tag "$tag" "Release $version" "$dry_run"; then
            success=false
            break
        fi
    done

    if [[ "$success" == "false" ]]; then
        if [[ "$dry_run" != "1" ]]; then
            log_error "Tag creation failed!"
            rollback_tags
        fi
        return 1
    fi

    if [[ "$dry_run" == "1" ]]; then
        log_info "[DRY RUN] All $(count_modules) tags would be created"
    else
        save_created_tags
        log_success "All $(count_modules) tags created locally"
    fi

    return 0
}

# Rollback - delete all created tags
rollback_tags() {
    log_section "Rolling Back Tags"

    load_created_tags

    if [[ ${#CREATED_TAGS[@]} -eq 0 ]]; then
        log_info "No tags to rollback"
        return 0
    fi

    log_warning "Deleting ${#CREATED_TAGS[@]} created tags..."

    for tag in "${CREATED_TAGS[@]}"; do
        if git tag -d "$tag" 2>/dev/null; then
            log_success "Deleted $tag"
        else
            log_warning "Could not delete $tag (may not exist)"
        fi
    done

    # Clean up tags file
    rm -f "$TAGS_FILE"

    log_warning "Rollback complete. No tags were pushed to remote."
}

# Push all tags to remote
push_tags() {
    local dry_run="${1:-0}"

    load_created_tags

    if [[ ${#CREATED_TAGS[@]} -eq 0 ]]; then
        log_error "No tags to push"
        return 1
    fi

    log_section "Pushing Tags to Remote"

    if [[ "$dry_run" == "1" ]]; then
        log_info "[DRY RUN] Would push ${#CREATED_TAGS[@]} tags to origin:"
        for tag in "${CREATED_TAGS[@]}"; do
            echo "    $tag"
        done
        return 0
    fi

    show_progress "Pushing ${#CREATED_TAGS[@]} tags to origin..."

    # Push all tags in a single operation
    if git push origin "${CREATED_TAGS[@]}" 2>&1 | tee /tmp/fuego-release-push.log; then
        log_success "All tags pushed successfully!"

        # Clean up tags file after successful push
        rm -f "$TAGS_FILE"

        echo ""
        log_info "View tags on GitHub:"
        for tag in "${CREATED_TAGS[@]}"; do
            echo "  https://github.com/go-fuego/fuego/releases/tag/$tag"
        done

        return 0
    else
        log_error "Failed to push tags"
        log_info "Tags remain local. You can:"
        log_info "  1. Fix the issue and retry: make release-push"
        log_info "  2. Rollback: make release-rollback"
        log_info "Push log saved to: /tmp/fuego-release-push.log"
        return 1
    fi
}

# Check if tags were already created
tags_exist_locally() {
    load_created_tags
    [[ ${#CREATED_TAGS[@]} -gt 0 ]]
}

# Verify all tags before creating
verify_tags_available() {
    local version="$1"

    log_section "Verifying Tags"

    local all_available=true

    for entry in "${RELEASABLE_MODULES[@]}"; do
        local module_path tag_prefix
        read -r module_path tag_prefix <<< "$(parse_module_entry "$entry")"

        local tag
        tag=$(get_module_tag "$tag_prefix" "$version")

        if git tag -l "$tag" | grep -q "^${tag}$"; then
            log_error "Tag already exists: $tag"
            all_available=false
        fi
    done

    if [[ "$all_available" == "true" ]]; then
        log_success "All tags available"
        return 0
    else
        log_error "Some tags already exist. Cannot proceed."
        return 1
    fi
}

# Main script
main() {
    local mode="${1:-}"
    local version="${2:-}"
    local dry_run=0

    case "$mode" in
        --create)
            shift
            version="$1"
            dry_run="${2:-0}"

            if [[ -z "$version" ]]; then
                die "Version required for --create"
            fi

            verify_tags_available "$version" || exit 1
            preview_tags "$version"

            if [[ "$dry_run" != "1" ]]; then
                echo ""
                if ! confirm "Create these tags?" "n"; then
                    log_info "Cancelled"
                    exit 0
                fi
            fi

            create_all_tags "$version" "$dry_run"
            ;;

        --push)
            dry_run="${2:-0}"

            if [[ "$dry_run" == "1" ]]; then
                push_tags 1
            else
                if ! tags_exist_locally; then
                    die "No tags found to push. Create tags first with --create"
                fi

                echo ""
                if confirm "Push ${#CREATED_TAGS[@]} tags to origin?" "n"; then
                    push_tags 0
                else
                    log_info "Cancelled. Tags remain local."
                    exit 0
                fi
            fi
            ;;

        --rollback)
            if ! tags_exist_locally; then
                log_info "No tags to rollback"
                exit 0
            fi

            echo ""
            if confirm "Delete ${#CREATED_TAGS[@]} local tags?" "n"; then
                rollback_tags
            else
                log_info "Cancelled"
                exit 0
            fi
            ;;

        --preview)
            shift
            version="$1"

            if [[ -z "$version" ]]; then
                die "Version required for --preview"
            fi

            preview_tags "$version"
            ;;

        --list-created)
            load_created_tags
            if [[ ${#CREATED_TAGS[@]} -eq 0 ]]; then
                echo "No tags created"
            else
                echo "Created tags (${#CREATED_TAGS[@]}):"
                for tag in "${CREATED_TAGS[@]}"; do
                    echo "  $tag"
                done
            fi
            ;;

        --help|-h)
            cat <<EOF
Usage: tag.sh [MODE] [VERSION] [OPTIONS]

Tag management for Fuego multi-module releases.

MODES:
    --create VERSION [--dry-run]  Create all tags for version
    --push [--dry-run]            Push created tags to remote
    --rollback                    Delete all created local tags
    --preview VERSION             Preview tags that would be created
    --list-created                List tags that were created
    --help, -h                    Show this help message

EXAMPLES:
    # Preview tags
    ./tag.sh --preview v0.19.0

    # Create tags (with confirmation)
    ./tag.sh --create v0.19.0

    # Dry run (no actual tags created)
    ./tag.sh --create v0.19.0 --dry-run

    # Push tags to remote
    ./tag.sh --push

    # Rollback (delete local tags)
    ./tag.sh --rollback

    # List created tags
    ./tag.sh --list-created
EOF
            ;;

        *)
            log_error "Unknown mode: $mode"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
}

# Run main if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi

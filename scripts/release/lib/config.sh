#!/bin/bash

# Configuration for multi-module releases
# Defines releasable modules and provides module utilities

# Prevent double-sourcing
if [[ -n "${_FUEGO_CONFIG_LOADED:-}" ]]; then
    return 0
fi
_FUEGO_CONFIG_LOADED=1

set -euo pipefail

# Source common utilities
_CONFIG_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./common.sh
source "$_CONFIG_LIB_DIR/common.sh"

# Define all releasable modules (9 modules)
# Format: "module_path:tag_prefix"
# - module_path: relative path from repo root ("." for root)
# - tag_prefix: prefix for git tag ("." for root means no prefix)
declare -a RELEASABLE_MODULES=(
    ".:."                                      # Main: v0.19.0
    "cmd/fuego:cmd/fuego"                      # CLI: cmd/fuego/v0.19.0
    "extra/fuegoecho:extra/fuegoecho"          # extra/fuegoecho/v0.19.0
    "extra/fuegogin:extra/fuegogin"            # extra/fuegogin/v0.19.0
    "extra/markdown:extra/markdown"            # extra/markdown/v0.19.0
    "extra/sql:extra/sql"                      # extra/sql/v0.19.0
    "extra/sqlite3:extra/sqlite3"              # extra/sqlite3/v0.19.0
    "middleware/basicauth:middleware/basicauth" # middleware/basicauth/v0.19.0
    "middleware/cache:middleware/cache"        # middleware/cache/v0.19.0
)

# Get git root directory
get_git_root_dir() {
    git rev-parse --show-toplevel
}

# Get absolute path for a module
get_module_dir() {
    local module_path="$1"
    local git_root
    git_root=$(get_git_root_dir)

    if [[ "$module_path" == "." ]]; then
        echo "$git_root"
    else
        echo "$git_root/$module_path"
    fi
}

# Get tag name for a module and version
get_module_tag() {
    local tag_prefix="$1"
    local version="$2"

    if [[ "$tag_prefix" == "." ]]; then
        echo "$version"
    else
        echo "$tag_prefix/$version"
    fi
}

# Parse module entry (module_path:tag_prefix)
# Returns: prints "module_path tag_prefix" separated by space
parse_module_entry() {
    local entry="$1"
    echo "$entry" | tr ':' ' '
}

# Iterate over all releasable modules
# Usage: for_each_module callback_function
for_each_module() {
    local callback="$1"
    shift  # Remove callback from args, rest are passed to callback

    for entry in "${RELEASABLE_MODULES[@]}"; do
        local module_path tag_prefix
        read -r module_path tag_prefix <<< "$(parse_module_entry "$entry")"

        # Call callback with module info and any additional args
        "$callback" "$module_path" "$tag_prefix" "$@"
    done
}

# Count releasable modules
count_modules() {
    echo "${#RELEASABLE_MODULES[@]}"
}

# Verify all modules exist and have go.mod
verify_modules_exist() {
    local all_exist=true
    local missing_modules=()

    for entry in "${RELEASABLE_MODULES[@]}"; do
        local module_path tag_prefix
        read -r module_path tag_prefix <<< "$(parse_module_entry "$entry")"

        local module_dir
        module_dir=$(get_module_dir "$module_path")

        if [[ ! -d "$module_dir" ]]; then
            all_exist=false
            missing_modules+=("$module_path (directory not found)")
        elif [[ ! -f "$module_dir/go.mod" ]]; then
            all_exist=false
            missing_modules+=("$module_path (go.mod not found)")
        fi
    done

    if [[ "$all_exist" == "false" ]]; then
        log_error "Some modules are missing or invalid:"
        for missing in "${missing_modules[@]}"; do
            log_error "  - $missing"
        done
        return 1
    fi

    return 0
}

# Get module name from go.mod
get_module_name() {
    local module_path="$1"
    local module_dir
    module_dir=$(get_module_dir "$module_path")

    if [[ ! -f "$module_dir/go.mod" ]]; then
        return 1
    fi

    # Extract module name from go.mod (first line: module github.com/...)
    grep -m1 "^module " "$module_dir/go.mod" | awk '{print $2}'
}

# Display all modules
list_modules() {
    log_section "Releasable Modules ($(count_modules) total)"

    local index=1
    for entry in "${RELEASABLE_MODULES[@]}"; do
        local module_path tag_prefix
        read -r module_path tag_prefix <<< "$(parse_module_entry "$entry")"

        local module_name
        module_name=$(get_module_name "$module_path" 2>/dev/null || echo "unknown")

        local display_path="$module_path"
        [[ "$display_path" == "." ]] && display_path="(root)"

        echo "  $index. $display_path"
        echo "     Module: $module_name"
        echo "     Tag pattern: $(get_module_tag "$tag_prefix" "vX.Y.Z")"
        echo ""
        index=$((index + 1))
    done
}

# Export functions for use in other scripts
export -f get_git_root_dir
export -f get_module_dir
export -f get_module_tag
export -f parse_module_entry
export -f for_each_module
export -f count_modules
export -f verify_modules_exist
export -f get_module_name
export -f list_modules

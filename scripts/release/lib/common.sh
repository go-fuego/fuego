#!/bin/bash

# Common utilities for release scripts
# Provides logging, colors, user interaction, and validation functions

# Prevent double-sourcing
if [[ -n "${_FUEGO_COMMON_LOADED:-}" ]]; then
    return 0
fi
_FUEGO_COMMON_LOADED=1

set -euo pipefail

# Colors for output
if [[ -z "${COLOR_RESET:-}" ]]; then
    readonly COLOR_RESET='\033[0m'
    readonly COLOR_RED='\033[0;31m'
    readonly COLOR_GREEN='\033[0;32m'
    readonly COLOR_YELLOW='\033[0;33m'
    readonly COLOR_BLUE='\033[0;34m'
    readonly COLOR_GRAY='\033[0;90m'
fi

# Logging functions
log_info() {
    echo -e "${COLOR_BLUE}ℹ${COLOR_RESET} $*"
}

log_success() {
    echo -e "${COLOR_GREEN}✓${COLOR_RESET} $*"
}

log_warning() {
    echo -e "${COLOR_YELLOW}⚠${COLOR_RESET} $*"
}

log_error() {
    echo -e "${COLOR_RED}✗${COLOR_RESET} $*" >&2
}

log_debug() {
    if [[ "${DEBUG:-0}" == "1" ]]; then
        echo -e "${COLOR_GRAY}[DEBUG]${COLOR_RESET} $*" >&2
    fi
}

log_section() {
    echo ""
    echo -e "${COLOR_BLUE}$*${COLOR_RESET}"
    printf '=%.0s' {1..60}
    echo ""
}

# Exit with error message
die() {
    log_error "$*"
    exit 1
}

# User interaction - confirm prompt (Y/n)
confirm() {
    local prompt="$1"
    local default="${2:-n}"

    if [[ "${FUEGO_RELEASE_AUTO_CONFIRM:-0}" == "1" ]]; then
        log_debug "Auto-confirm enabled, skipping prompt: $prompt"
        return 0
    fi

    local yn_prompt
    if [[ "$default" == "y" ]]; then
        yn_prompt="[Y/n]"
    else
        yn_prompt="[y/N]"
    fi

    while true; do
        read -rp "$(echo -e "${COLOR_YELLOW}?${COLOR_RESET} $prompt $yn_prompt ") " response
        response="${response:-$default}"
        case "$response" in
            [Yy]* ) return 0;;
            [Nn]* ) return 1;;
            * ) echo "Please answer yes or no.";;
        esac
    done
}

# Get user input
prompt() {
    local prompt_text="$1"
    local default="${2:-}"

    if [[ -n "$default" ]]; then
        read -rp "$(echo -e "${COLOR_YELLOW}?${COLOR_RESET} $prompt_text [$default]: ") " response
        echo "${response:-$default}"
    else
        read -rp "$(echo -e "${COLOR_YELLOW}?${COLOR_RESET} $prompt_text: ") " response
        echo "$response"
    fi
}

# Validation functions
is_semantic_version() {
    local version="$1"
    # Match vX.Y.Z format (semantic versioning)
    if [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        return 0
    else
        return 1
    fi
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Git helpers
get_git_root() {
    git rev-parse --show-toplevel 2>/dev/null
}

is_git_clean() {
    [[ -z "$(git status --porcelain)" ]]
}

has_unpushed_commits() {
    local branch="${1:-main}"
    # Check if branch has commits not in origin/branch
    local unpushed
    unpushed=$(git log "origin/$branch..$branch" --oneline 2>/dev/null | wc -l)
    [[ "$unpushed" -gt 0 ]]
}

get_current_branch() {
    git rev-parse --abbrev-ref HEAD
}

is_branch_up_to_date() {
    local branch="${1:-main}"
    git fetch origin "$branch" --quiet 2>/dev/null || return 1

    local local_commit
    local remote_commit
    local_commit=$(git rev-parse "$branch")
    remote_commit=$(git rev-parse "origin/$branch")

    [[ "$local_commit" == "$remote_commit" ]]
}

# Version comparison helpers
version_compare() {
    local v1="$1"
    local v2="$2"

    # Remove 'v' prefix
    v1="${v1#v}"
    v2="${v2#v}"

    # Split into components
    local v1_major v1_minor v1_patch
    local v2_major v2_minor v2_patch
    IFS='.' read -r v1_major v1_minor v1_patch <<< "$v1"
    IFS='.' read -r v2_major v2_minor v2_patch <<< "$v2"

    # Remove leading zeros to avoid octal interpretation
    v1_major=$((10#${v1_major}))
    v1_minor=$((10#${v1_minor}))
    v1_patch=$((10#${v1_patch}))
    v2_major=$((10#${v2_major}))
    v2_minor=$((10#${v2_minor}))
    v2_patch=$((10#${v2_patch}))

    # Compare major version
    if [[ "$v1_major" -gt "$v2_major" ]]; then
        echo "gt"
        return
    elif [[ "$v1_major" -lt "$v2_major" ]]; then
        echo "lt"
        return
    fi

    # Major versions equal, compare minor
    if [[ "$v1_minor" -gt "$v2_minor" ]]; then
        echo "gt"
        return
    elif [[ "$v1_minor" -lt "$v2_minor" ]]; then
        echo "lt"
        return
    fi

    # Major and minor equal, compare patch
    if [[ "$v1_patch" -gt "$v2_patch" ]]; then
        echo "gt"
    elif [[ "$v1_patch" -lt "$v2_patch" ]]; then
        echo "lt"
    else
        echo "eq"
    fi
}

# Increment version
increment_version() {
    local version="$1"
    local part="${2:-patch}"  # patch, minor, or major

    # Remove 'v' prefix
    version="${version#v}"

    local major minor patch
    IFS='.' read -r major minor patch <<< "$version"

    case "$part" in
        patch)
            patch=$((patch + 1))
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        *)
            die "Invalid version part: $part (must be patch, minor, or major)"
            ;;
    esac

    echo "v${major}.${minor}.${patch}"
}

# Check if running from git root
ensure_git_root() {
    local git_root
    git_root=$(get_git_root) || die "Not in a git repository"

    local current_dir
    current_dir=$(pwd)

    if [[ "$current_dir" != "$git_root" ]]; then
        die "Must run from git repository root. Current: $current_dir, Root: $git_root"
    fi
}

# Print banner
print_banner() {
    local title="$1"
    echo ""
    echo -e "${COLOR_BLUE}🚀 $title${COLOR_RESET}"
    printf '=%.0s' {1..60}
    echo ""
}

# Progress indicator
show_progress() {
    local message="$1"
    echo -e "${COLOR_BLUE}⏳${COLOR_RESET} $message"
}

#!/bin/bash

# Pre-release validation for multi-module releases
# Runs comprehensive checks before allowing a release

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"
# shellcheck source=./lib/config.sh
source "$SCRIPT_DIR/lib/config.sh"

# Track validation status
VALIDATION_FAILED=0

# Check if command exists and show version
check_command() {
    local cmd="$1"
    local display_name="${2:-$cmd}"

    if command_exists "$cmd"; then
        local version
        version=$("$cmd" --version 2>&1 | head -n1 || echo "unknown")
        log_success "$display_name is installed ($version)"
        return 0
    else
        log_error "$display_name is not installed"
        VALIDATION_FAILED=1
        return 1
    fi
}

# Environment checks
validate_environment() {
    log_section "Environment Checks"

    check_command "git" "Git"
    check_command "go" "Go"
    check_command "golangci-lint" "golangci-lint"

    # Check if in git repo
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        log_error "Not in a git repository"
        VALIDATION_FAILED=1
        return 1
    fi
    log_success "In git repository"

    # Check if at repo root
    local git_root current_dir
    git_root=$(get_git_root)
    current_dir=$(pwd)

    if [[ "$current_dir" != "$git_root" ]]; then
        log_error "Not at repository root"
        log_info "Current: $current_dir"
        log_info "Root: $git_root"
        log_info "Run: cd $git_root"
        VALIDATION_FAILED=1
        return 1
    fi
    log_success "At repository root"

    # Verify all modules exist
    if verify_modules_exist; then
        log_success "All $(count_modules) modules found"
    else
        VALIDATION_FAILED=1
        return 1
    fi

    return 0
}

# Git state checks
validate_git_state() {
    log_section "Git State Checks"

    # Check working directory is clean
    if is_git_clean; then
        log_success "Working directory is clean"
    else
        log_error "Working directory has uncommitted changes"
        log_info "Please commit or stash your changes:"
        log_info "  git status"
        VALIDATION_FAILED=1
        return 1
    fi

    # Check current branch
    local current_branch
    current_branch=$(get_current_branch)

    if [[ "$current_branch" == "main" ]]; then
        log_success "On main branch"
    else
        log_warning "Not on main branch (current: $current_branch)"
        if [[ "${FUEGO_RELEASE_ALLOW_NON_MAIN:-0}" != "1" ]]; then
            log_error "Release must be done from main branch"
            log_info "To override: export FUEGO_RELEASE_ALLOW_NON_MAIN=1"
            VALIDATION_FAILED=1
            return 1
        fi
    fi

    # Fetch latest from origin
    show_progress "Fetching latest from origin..."
    if ! git fetch origin main --quiet 2>/dev/null; then
        log_warning "Could not fetch from origin (network issue?)"
    fi

    # Check if up to date with origin/main
    if is_branch_up_to_date "main"; then
        log_success "Up-to-date with origin/main"
    else
        log_error "Not up-to-date with origin/main"
        log_info "Please pull latest changes:"
        log_info "  git pull origin main"
        VALIDATION_FAILED=1
        return 1
    fi

    # Check for unpushed commits
    if has_unpushed_commits "main"; then
        log_error "You have unpushed commits"
        log_info "Please push your commits before releasing:"
        log_info "  git push origin main"
        VALIDATION_FAILED=1
        return 1
    fi
    log_success "No unpushed commits"

    return 0
}

# Validate version
validate_version_arg() {
    local version="$1"

    log_section "Version Validation"

    # Run version.sh validation
    if "$SCRIPT_DIR/version.sh" --validate "$version" >/dev/null 2>&1; then
        log_success "Version $version is valid and available"
        return 0
    else
        log_error "Version $version is invalid or unavailable"
        VALIDATION_FAILED=1
        return 1
    fi
}

# Run linting
run_lint() {
    log_section "Running Linter"

    show_progress "Running golangci-lint (this may take a minute)..."

    local git_root
    git_root=$(get_git_root)

    if make -C "$git_root" lint 2>&1 | tee /tmp/fuego-release-lint.log; then
        log_success "Lint passed"
        return 0
    else
        log_error "Lint failed"
        log_info "Check output above or see: /tmp/fuego-release-lint.log"
        log_info "Fix issues and try again"
        VALIDATION_FAILED=1
        return 1
    fi
}

# Run tests
run_tests() {
    log_section "Running Tests"

    show_progress "Running tests for all modules..."

    local git_root
    git_root=$(get_git_root)

    # Use existing test target with full paths
    if make -C "$git_root" test 2>&1 | tee /tmp/fuego-release-test.log; then
        log_success "All tests passed"
        return 0
    else
        log_error "Tests failed"
        log_info "Check output above or see: /tmp/fuego-release-test.log"
        log_info "Fix failing tests and try again"
        VALIDATION_FAILED=1
        return 1
    fi
}

# Run builds
run_build() {
    log_section "Running Builds"

    show_progress "Building all modules..."

    local git_root
    git_root=$(get_git_root)

    if make -C "$git_root" build 2>&1 | tee /tmp/fuego-release-build.log; then
        log_success "All builds succeeded"
        return 0
    else
        log_error "Build failed"
        log_info "Check output above or see: /tmp/fuego-release-build.log"
        log_info "Fix build errors and try again"
        VALIDATION_FAILED=1
        return 1
    fi
}

# Comprehensive validation
run_comprehensive_validation() {
    local skip_lint="${FUEGO_RELEASE_SKIP_LINT:-0}"
    local skip_tests="${FUEGO_RELEASE_SKIP_TESTS:-0}"
    local skip_build="${FUEGO_RELEASE_SKIP_BUILD:-0}"

    if [[ "$skip_lint" != "1" ]]; then
        run_lint || true
    else
        log_warning "Skipping lint (FUEGO_RELEASE_SKIP_LINT=1)"
    fi

    if [[ "$skip_tests" != "1" ]]; then
        run_tests || true
    else
        log_warning "Skipping tests (FUEGO_RELEASE_SKIP_TESTS=1)"
    fi

    if [[ "$skip_build" != "1" ]]; then
        run_build || true
    else
        log_warning "Skipping build (FUEGO_RELEASE_SKIP_BUILD=1)"
    fi
}

# Main validation flow
main() {
    local mode="${1:-full}"
    local version="${2:-}"

    print_banner "Fuego Pre-Release Validation"

    case "$mode" in
        --env-only)
            validate_environment
            ;;
        --git-only)
            validate_environment
            validate_git_state
            ;;
        --validate-only)
            # Skip git checks, only run lint/test/build
            run_comprehensive_validation
            ;;
        --version)
            shift
            version="$1"
            validate_environment
            validate_version_arg "$version"
            ;;
        --help|-h)
            cat <<EOF
Usage: validate.sh [MODE] [VERSION]

Pre-release validation for Fuego multi-module releases.

MODES:
    (no mode)           Full validation (default)
    --env-only          Only environment checks
    --git-only          Environment + git state checks
    --validate-only     Only run lint/test/build (skip git checks)
    --version VERSION   Validate a specific version
    --help, -h          Show this help message

ENVIRONMENT VARIABLES:
    FUEGO_RELEASE_SKIP_LINT=1       Skip linting
    FUEGO_RELEASE_SKIP_TESTS=1      Skip tests
    FUEGO_RELEASE_SKIP_BUILD=1      Skip builds
    FUEGO_RELEASE_ALLOW_NON_MAIN=1  Allow release from non-main branch

EXAMPLES:
    # Full validation
    ./validate.sh

    # Check git state only
    ./validate.sh --git-only

    # Run tests and builds only
    ./validate.sh --validate-only

    # Validate specific version
    ./validate.sh --version v0.19.0
EOF
            exit 0
            ;;
        full|*)
            # Full validation flow
            validate_environment
            validate_git_state

            if [[ -n "$version" ]]; then
                validate_version_arg "$version"
            fi

            run_comprehensive_validation
            ;;
    esac

    echo ""
    if [[ "$VALIDATION_FAILED" -eq 0 ]]; then
        log_section "Validation Complete"
        log_success "All pre-release checks passed!"
        echo ""
        exit 0
    else
        log_section "Validation Failed"
        log_error "Some checks failed. Please fix the issues above and try again."
        echo ""
        exit 2
    fi
}

# Run main if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi

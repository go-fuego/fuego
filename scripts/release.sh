#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
DRY_RUN=false
FORCE=false
SKIP_TESTS=false

# Function to print colored output
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_header() {
    echo -e "${BLUE}🚀 $1${NC}"
}

# Function to show usage
show_usage() {
    cat << EOF
Usage: $0 [OPTIONS] VERSION

Create a release with tags for all modules in the workspace.

ARGUMENTS:
    VERSION     Release version (e.g., v1.2.3)

OPTIONS:
    -d, --dry-run       Preview changes without actually creating tags
    -f, --force         Skip confirmation prompts
    -s, --skip-tests    Skip running tests before release
    -h, --help          Show this help message

EXAMPLES:
    $0 v1.2.3                    # Create release v1.2.3
    $0 --dry-run v1.2.3          # Preview release v1.2.3
    $0 --force v1.2.3            # Create release without confirmation
    $0 --skip-tests v1.2.3       # Create release without running tests

MODULES:
    This script will automatically discover and tag all modules in go.work:
    - Main module: v1.2.3
    - Extra modules: extra/*/v1.2.3

EOF
}

# Function to validate version format
validate_version() {
    local version="$1"
    if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]]; then
        print_error "Invalid version format: $version"
        print_info "Expected format: v1.2.3 or v1.2.3-alpha.1"
        exit 1
    fi
}

# Function to check if git working directory is clean
check_git_status() {
    if [[ -n $(git status --porcelain) ]]; then
        print_error "Working directory is not clean. Please commit or stash your changes."
        git status --short
        exit 1
    fi
}

# Function to check if tags already exist
check_existing_tags() {
    local version="$1"
    local existing_tags=()
    
    # Check main module tag
    if git tag -l | grep -q "^${version}$"; then
        existing_tags+=("$version")
    fi
    
    # Check extra module tags
    for module in extra/fuegoecho extra/fuegogin extra/markdown extra/sql extra/sqlite3; do
        local tag="${module}/${version}"
        if git tag -l | grep -q "^${tag}$"; then
            existing_tags+=("$tag")
        fi
    done
    
    if [[ ${#existing_tags[@]} -gt 0 ]]; then
        print_error "The following tags already exist:"
        for tag in "${existing_tags[@]}"; do
            echo "  - $tag"
        done
        exit 1
    fi
}

# Function to discover modules from go.work
discover_modules() {
    if [[ ! -f "$PROJECT_ROOT/go.work" ]]; then
        print_error "go.work file not found in project root"
        exit 1
    fi
    
    cd "$PROJECT_ROOT"
    local modules
    modules=$(grep -E "^\s*\./|^\s*\." go.work | grep -E "extra/" | sed 's/^\s*\.\///' | tr '\n' ' ')
    echo "$modules"
}

# Function to run tests
run_tests() {
    if [[ "$SKIP_TESTS" == "true" ]]; then
        print_warning "Skipping tests as requested"
        return
    fi
    
    print_info "Running tests..."
    cd "$PROJECT_ROOT"
    
    if ! make ci; then
        print_error "CI tests failed"
        exit 1
    fi
    
    print_success "All tests passed"
}

# Function to generate changelog
generate_changelog() {
    local version="$1"
    local changelog_file="$PROJECT_ROOT/CHANGELOG_${version}.md"
    
    print_info "Generating changelog..."
    
    # Get the last release tag
    local last_tag
    last_tag=$(git tag -l "v*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+$" | sort -V | tail -n 1)
    
    if [[ -z "$last_tag" ]]; then
        last_tag=$(git rev-list --max-parents=0 HEAD)
        print_info "No previous release found, using first commit"
    else
        print_info "Last release: $last_tag"
    fi
    
    # Generate changelog
    {
        echo "# Release $version"
        echo ""
        echo "## What's Changed"
        echo ""
        git log --pretty=format:"* %s (%h)" "${last_tag}..HEAD"
        echo ""
        echo ""
        echo "## Contributors"
        echo ""
        git log --pretty=format:"* @%an" "${last_tag}..HEAD" | sort -u
    } > "$changelog_file"
    
    print_success "Changelog generated: $changelog_file"
    
    if [[ "$DRY_RUN" == "false" ]]; then
        print_info "Preview of changelog:"
        echo "----------------------------------------"
        head -20 "$changelog_file"
        echo "----------------------------------------"
    fi
}

# Function to create tags
create_tags() {
    local version="$1"
    local modules="$2"
    
    print_header "Creating tags for version: $version"
    
    # Create main module tag
    print_info "Creating main tag: $version"
    if [[ "$DRY_RUN" == "true" ]]; then
        echo "[DRY RUN] Would create: git tag -a $version -m 'Release $version'"
    else
        git tag -a "$version" -m "Release $version"
        print_success "Created tag: $version"
    fi
    
    # Create extra module tags
    for module in $modules; do
        local tag="${module}/${version}"
        print_info "Creating module tag: $tag"
        if [[ "$DRY_RUN" == "true" ]]; then
            echo "[DRY RUN] Would create: git tag -a $tag -m 'Release $tag'"
        else
            git tag -a "$tag" -m "Release $tag"
            print_success "Created tag: $tag"
        fi
    done
}

# Function to push tags
push_tags() {
    if [[ "$DRY_RUN" == "true" ]]; then
        print_info "[DRY RUN] Would push all tags to origin"
        return
    fi
    
    print_info "Pushing tags to origin..."
    git push origin --tags
    print_success "All tags pushed to origin"
}

# Function to show summary
show_summary() {
    local version="$1"
    local modules="$2"
    
    print_header "Release Summary"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_warning "DRY RUN COMPLETED - No actual changes were made"
        echo ""
        echo "The following would have been created:"
    else
        print_success "RELEASE COMPLETED"
        echo ""
        echo "Successfully created the following releases:"
    fi
    
    echo ""
    echo "🏷️  Tags Created:"
    echo "  - $version (main module)"
    
    for module in $modules; do
        echo "  - ${module}/${version} ($module module)"
    done
    
    if [[ "$DRY_RUN" == "false" ]]; then
        echo ""
        echo "📝 Next steps:"
        echo "  1. Check the releases on GitHub"
        echo "  2. Update documentation if needed"
        echo "  3. Announce the release"
    fi
}

# Function to confirm action
confirm_action() {
    local version="$1"
    local modules="$2"
    
    if [[ "$FORCE" == "true" ]]; then
        return
    fi
    
    echo ""
    print_header "Release Confirmation"
    echo "Version: $version"
    echo "Modules to be tagged:"
    echo "  - $version (main)"
    for module in $modules; do
        echo "  - ${module}/${version}"
    done
    echo ""
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_warning "This is a DRY RUN - no actual changes will be made"
    else
        print_warning "This will create and push tags to the remote repository"
    fi
    
    echo ""
    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Release cancelled"
        exit 0
    fi
}

# Function to cleanup on error
cleanup_on_error() {
    local version="$1"
    local modules="$2"
    
    print_error "An error occurred during the release process"
    print_info "Cleaning up any created tags..."
    
    # Remove main tag if it exists
    if git tag -l | grep -q "^${version}$"; then
        git tag -d "$version" 2>/dev/null || true
        print_info "Removed tag: $version"
    fi
    
    # Remove module tags if they exist
    for module in $modules; do
        local tag="${module}/${version}"
        if git tag -l | grep -q "^${tag}$"; then
            git tag -d "$tag" 2>/dev/null || true
            print_info "Removed tag: $tag"
        fi
    done
    
    print_info "Cleanup completed"
    exit 1
}

# Main function
main() {
    local version=""
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -f|--force)
                FORCE=true
                shift
                ;;
            -s|--skip-tests)
                SKIP_TESTS=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            -*)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
            *)
                if [[ -z "$version" ]]; then
                    version="$1"
                else
                    print_error "Too many arguments"
                    show_usage
                    exit 1
                fi
                shift
                ;;
        esac
    done
    
    # Check if version is provided
    if [[ -z "$version" ]]; then
        print_error "Version is required"
        show_usage
        exit 1
    fi
    
    # Change to project root
    cd "$PROJECT_ROOT"
    
    # Validate inputs
    validate_version "$version"
    
    # Check git status
    if [[ "$DRY_RUN" == "false" ]]; then
        check_git_status
    fi
    
    # Check existing tags
    check_existing_tags "$version"
    
    # Discover modules
    local modules
    modules=$(discover_modules)
    print_info "Discovered modules: $modules"
    
    # Confirm action
    confirm_action "$version" "$modules"
    
    # Set up error handling
    if [[ "$DRY_RUN" == "false" ]]; then
        trap 'cleanup_on_error "$version" "$modules"' ERR
    fi
    
    # Run tests
    run_tests
    
    # Generate changelog
    generate_changelog "$version"
    
    # Create tags
    create_tags "$version" "$modules"
    
    # Push tags
    push_tags
    
    # Show summary
    show_summary "$version" "$modules"
    
    # Disable error trap
    trap - ERR
}

# Run main function
main "$@"

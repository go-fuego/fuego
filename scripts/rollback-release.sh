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
    echo -e "${BLUE}🔄 $1${NC}"
}

# Function to show usage
show_usage() {
    cat << EOF
Usage: $0 [OPTIONS] VERSION

Rollback a release by deleting tags for all modules in the workspace.

ARGUMENTS:
    VERSION     Release version to rollback (e.g., v1.2.3)

OPTIONS:
    -d, --dry-run       Preview changes without actually deleting tags
    -f, --force         Skip confirmation prompts
    -h, --help          Show this help message

EXAMPLES:
    $0 v1.2.3                    # Rollback release v1.2.3
    $0 --dry-run v1.2.3          # Preview rollback of v1.2.3
    $0 --force v1.2.3            # Rollback without confirmation

WARNING:
    This script will delete both local and remote tags. Use with caution!
    GitHub releases will NOT be automatically deleted - you'll need to do that manually.

MODULES:
    This script will attempt to delete tags for all modules:
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

# Function to check which tags exist
check_existing_tags() {
    local version="$1"
    local modules="$2"
    local existing_local=()
    local existing_remote=()
    
    print_info "Checking which tags exist..."
    
    # Check main module tag
    if git tag -l | grep -q "^${version}$"; then
        existing_local+=("$version")
    fi
    if git ls-remote --tags origin | grep -q "refs/tags/${version}$"; then
        existing_remote+=("$version")
    fi
    
    # Check extra module tags
    for module in $modules; do
        local tag="${module}/${version}"
        if git tag -l | grep -q "^${tag}$"; then
            existing_local+=("$tag")
        fi
        if git ls-remote --tags origin | grep -q "refs/tags/${tag}$"; then
            existing_remote+=("$tag")
        fi
    done
    
    if [[ ${#existing_local[@]} -eq 0 && ${#existing_remote[@]} -eq 0 ]]; then
        print_warning "No tags found for version $version"
        print_info "Nothing to rollback"
        exit 0
    fi
    
    if [[ ${#existing_local[@]} -gt 0 ]]; then
        print_info "Local tags found:"
        for tag in "${existing_local[@]}"; do
            echo "  - $tag"
        done
    fi
    
    if [[ ${#existing_remote[@]} -gt 0 ]]; then
        print_info "Remote tags found:"
        for tag in "${existing_remote[@]}"; do
            echo "  - $tag"
        done
    fi
    
    # Return arrays via global variables (bash limitation)
    EXISTING_LOCAL=("${existing_local[@]}")
    EXISTING_REMOTE=("${existing_remote[@]}")
}

# Function to confirm action
confirm_action() {
    local version="$1"
    
    if [[ "$FORCE" == "true" ]]; then
        return
    fi
    
    echo ""
    print_header "Rollback Confirmation"
    echo "Version: $version"
    echo ""
    
    if [[ ${#EXISTING_LOCAL[@]} -gt 0 ]]; then
        echo "Local tags to be deleted:"
        for tag in "${EXISTING_LOCAL[@]}"; do
            echo "  - $tag"
        done
    fi
    
    if [[ ${#EXISTING_REMOTE[@]} -gt 0 ]]; then
        echo "Remote tags to be deleted:"
        for tag in "${EXISTING_REMOTE[@]}"; do
            echo "  - $tag"
        done
    fi
    
    echo ""
    if [[ "$DRY_RUN" == "true" ]]; then
        print_warning "This is a DRY RUN - no actual changes will be made"
    else
        print_warning "This will permanently delete the tags from local and remote repositories"
        print_warning "GitHub releases will NOT be deleted automatically"
    fi
    
    echo ""
    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Rollback cancelled"
        exit 0
    fi
}

# Function to delete local tags
delete_local_tags() {
    if [[ ${#EXISTING_LOCAL[@]} -eq 0 ]]; then
        print_info "No local tags to delete"
        return
    fi
    
    print_header "Deleting local tags"
    
    for tag in "${EXISTING_LOCAL[@]}"; do
        print_info "Deleting local tag: $tag"
        if [[ "$DRY_RUN" == "true" ]]; then
            echo "[DRY RUN] Would delete: git tag -d $tag"
        else
            if git tag -d "$tag"; then
                print_success "Deleted local tag: $tag"
            else
                print_error "Failed to delete local tag: $tag"
            fi
        fi
    done
}

# Function to delete remote tags
delete_remote_tags() {
    if [[ ${#EXISTING_REMOTE[@]} -eq 0 ]]; then
        print_info "No remote tags to delete"
        return
    fi
    
    print_header "Deleting remote tags"
    
    for tag in "${EXISTING_REMOTE[@]}"; do
        print_info "Deleting remote tag: $tag"
        if [[ "$DRY_RUN" == "true" ]]; then
            echo "[DRY RUN] Would delete: git push origin --delete $tag"
        else
            if git push origin --delete "$tag"; then
                print_success "Deleted remote tag: $tag"
            else
                print_error "Failed to delete remote tag: $tag"
            fi
        fi
    done
}

# Function to show summary
show_summary() {
    local version="$1"
    
    print_header "Rollback Summary"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_warning "DRY RUN COMPLETED - No actual changes were made"
        echo ""
        echo "The following would have been deleted:"
    else
        print_success "ROLLBACK COMPLETED"
        echo ""
        echo "Successfully deleted the following tags:"
    fi
    
    echo ""
    if [[ ${#EXISTING_LOCAL[@]} -gt 0 ]]; then
        echo "🏷️  Local tags:"
        for tag in "${EXISTING_LOCAL[@]}"; do
            echo "  - $tag"
        done
    fi
    
    if [[ ${#EXISTING_REMOTE[@]} -gt 0 ]]; then
        echo "🌐 Remote tags:"
        for tag in "${EXISTING_REMOTE[@]}"; do
            echo "  - $tag"
        done
    fi
    
    if [[ "$DRY_RUN" == "false" ]]; then
        echo ""
        print_warning "Manual steps still required:"
        echo "  1. Delete GitHub releases manually if needed"
        echo "  2. Update any documentation that references this version"
        echo "  3. Notify team members about the rollback"
    fi
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
    
    # Discover modules
    local modules
    modules=$(discover_modules)
    print_info "Discovered modules: $modules"
    
    # Check existing tags
    check_existing_tags "$version" "$modules"
    
    # Confirm action
    confirm_action "$version"
    
    # Delete tags
    delete_local_tags
    delete_remote_tags
    
    # Show summary
    show_summary "$version"
}

# Run main function
main "$@"

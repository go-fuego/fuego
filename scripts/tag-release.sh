#!/bin/bash

# Simple multi-module release tagger
# Usage: ./scripts/tag-release.sh v0.19.0

set -euo pipefail

VERSION="${1:-}"

if [[ -z "$VERSION" ]]; then
    echo "Usage: $0 vX.Y.Z"
    exit 1
fi

if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Invalid version format. Use vX.Y.Z (e.g., v0.19.0)"
    exit 1
fi

# Define all modules to tag
MODULES=(
    "."                        # Main: v0.19.0
    "cmd/fuego"                # cmd/fuego/v0.19.0
    "extra/fuegoecho"
    "extra/fuegogin"
    "extra/fuegomux"
    "extra/markdown"
    "extra/sql"
    "extra/sqlite3"
    "middleware/basicauth"
    "middleware/cache"
)

# Build tag list
tags=()
for module in "${MODULES[@]}"; do
    if [[ "$module" == "." ]]; then
        tags+=("$VERSION")
    else
        tags+=("$module/$VERSION")
    fi
done

echo "Tags to create for $VERSION:"
for tag in "${tags[@]}"; do
    echo "  $tag"
done
echo ""

# Check none already exist locally or remotely
echo "Checking for existing tags..."
conflicts=()
for tag in "${tags[@]}"; do
    if git rev-parse "refs/tags/$tag" &>/dev/null; then
        conflicts+=("$tag (local)")
    elif git ls-remote --tags origin "refs/tags/$tag" | grep -q .; then
        conflicts+=("$tag (remote)")
    fi
done

if [[ ${#conflicts[@]} -gt 0 ]]; then
    echo "Error: the following tags already exist:"
    for c in "${conflicts[@]}"; do
        echo "  $c"
    done
    exit 1
fi

# Create annotated tags
for tag in "${tags[@]}"; do
    git tag -a "$tag" -m "Release $VERSION"
    echo "  created $tag"
done

echo ""
read -r -p "Push all tags to origin? [y/N] " response
if [[ "$response" =~ ^[Yy]$ ]]; then
    for tag in "${tags[@]}"; do
        git push origin "$tag"
    done
    echo ""
    echo "Done. View at: https://github.com/go-fuego/fuego/tags"
else
    echo "Tags created locally only. Push with: git push origin --tags"
fi

#!/bin/bash

# Simple multi-module release tagger
# Usage: ./scripts/tag-release.sh v0.19.0

set -euo pipefail

VERSION="$1"

# Validate version format
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
    "extra/markdown"
    "extra/sql"
    "extra/sqlite3"
    "middleware/basicauth"
    "middleware/cache"
)

echo "Creating release $VERSION for all modules..."
echo ""

# Create all tags
for module in "${MODULES[@]}"; do
    if [[ "$module" == "." ]]; then
        tag="$VERSION"
    else
        tag="$module/$VERSION"
    fi

    echo "Creating tag: $tag"
    git tag -a "$tag" -m "Release $VERSION"
done

echo ""
echo "✓ All tags created locally"
echo ""
echo "Push to remote? [y/N] "
read -r response

if [[ "$response" =~ ^[Yy]$ ]]; then
    # Push all tags
    for module in "${MODULES[@]}"; do
        if [[ "$module" == "." ]]; then
            tag="$VERSION"
        else
            tag="$module/$VERSION"
        fi
        git push origin "$tag"
    done

    echo ""
    echo "✓ All tags pushed to origin"
    echo ""
    echo "View releases: https://github.com/go-fuego/fuego/tags"
else
    echo ""
    echo "Tags created locally but not pushed."
    echo "To push later: git push origin --tags"
fi

# Release Automation Guide

This document describes the automated release system for the Fuego project, which handles creating tags for all modules in the workspace.

## Overview

The Fuego project uses a Go workspace with multiple modules:

- **Main module**: `github.com/go-fuego/fuego`
- **Extra modules**:
  - `extra/fuegoecho` - Echo adapter
  - `extra/fuegogin` - Gin adapter
  - `extra/markdown` - Markdown utilities
  - `extra/sql` - SQL utilities
  - `extra/sqlite3` - SQLite3 utilities

Previously, creating a release required manually creating multiple git tags:

```bash
git tag -a v0.18.9
git tag -a extra/fuegogin/v0.18.9
git tag -a extra/fuegoecho/v0.18.9
# ... and so on for each module
```

The new automated system handles all of this with a single command or GitHub Actions workflow.

## Release Methods

### 1. GitHub Actions (Recommended for Team Releases)

**Location**: `.github/workflows/release.yml`

**How to use**:

1. Go to the **Actions** tab in your GitHub repository
2. Select the **"Release"** workflow
3. Click **"Run workflow"**
4. Enter the version (e.g., `v1.2.3`)
5. Optionally check "Dry run" to preview changes
6. Click **"Run workflow"**

**What it does**:

1. ✅ **Validates** version format and checks for existing tags
2. 🧪 **Runs tests** using `make ci` and `make ci-full`
3. ⏳ **Waits for manual approval** (unless dry run)
4. 🏷️ **Creates tags** for all modules
5. 📤 **Pushes tags** to GitHub
6. 📝 **Generates changelog** from git commits
7. 🚀 **Creates GitHub releases** with changelogs

**Features**:

- **Dry run mode**: Preview changes without creating actual releases
- **Manual approval**: Pause before creating tags for confirmation
- **Automatic module discovery**: Finds all modules from `go.work`
- **Comprehensive testing**: Runs full CI suite before release
- **Rich changelog**: Includes commits and contributors since last release
- **Multiple releases**: Creates separate GitHub releases for each module

### 2. Local Script (For Development and Testing)

**Location**: `scripts/release.sh`

**How to use**:

```bash
# Basic usage
./scripts/release.sh v1.2.3

# Preview changes (dry run)
./scripts/release.sh --dry-run v1.2.3

# Skip confirmation prompts
./scripts/release.sh --force v1.2.3

# Skip tests (not recommended)
./scripts/release.sh --skip-tests v1.2.3

# Get help
./scripts/release.sh --help
```

**Features**:

- 🎨 **Colored output** for better readability
- 🔍 **Dry run mode** to preview changes
- ✅ **Safety checks**: Clean working directory, existing tags
- 🧪 **Test integration**: Runs `make ci` before release
- 📝 **Changelog generation**: Creates `CHANGELOG_v1.2.3.md`
- 🛡️ **Error handling**: Automatic cleanup on failure
- 💬 **Interactive confirmations** (unless `--force`)

## Version Format

Both methods support semantic versioning:

- **Release versions**: `v1.2.3`, `v2.0.0`
- **Pre-release versions**: `v1.2.3-alpha.1`, `v1.2.3-beta.2`, `v1.2.3-rc.1`

## Tags Created

For version `v1.2.3`, the following tags are created:

- `v1.2.3` (main module)
- `extra/fuegoecho/v1.2.3`
- `extra/fuegogin/v1.2.3`
- `extra/markdown/v1.2.3`
- `extra/sql/v1.2.3`
- `extra/sqlite3/v1.2.3`

## GitHub Releases

Each module gets its own GitHub release:

- **Main release**: `v1.2.3` (marked as "latest")
- **Module releases**: `extra/fuegoecho/v1.2.3`, etc.

All releases include:

- 📝 **Changelog** with commits since last release
- 👥 **Contributors** list
- 🔗 **Links** to related releases

## Safety Features

### Pre-flight Checks

- ✅ Version format validation
- ✅ Clean working directory (local script only)
- ✅ No existing tags with same version
- ✅ All tests must pass

### Error Handling

- 🛡️ **Automatic cleanup**: Removes created tags on failure
- ⚠️ **Clear error messages**: Explains what went wrong
- 🔄 **Rollback capability**: Easy to undo partial releases

### Manual Approval

- ⏳ **GitHub Actions**: Pauses for manual approval before pushing tags
- 💬 **Local script**: Interactive confirmation (unless `--force`)
- 🔍 **Preview mode**: Dry run shows exactly what would be created

## Permissions

### GitHub Actions

- Requires `contents: write` permission (already configured)
- Any team member with repository access can trigger releases
- Consider setting up branch protection rules for additional security

### Local Script

- Requires git push access to the repository
- Must be run from the project root directory

## Troubleshooting

### Common Issues

**"Tag already exists"**

```
❌ Tag v1.2.3 already exists
```

**Solution**: Use a different version number or delete the existing tag if it was created in error.

**"Working directory is not clean"**

```
❌ Working directory is not clean. Please commit or stash your changes.
```

**Solution**: Commit or stash your changes before creating a release.

**"Tests failed"**

```
❌ CI tests failed
```

**Solution**: Fix the failing tests before creating a release. Use `--skip-tests` only for testing the release process itself.

### Manual Cleanup

If something goes wrong and you need to clean up tags manually:

```bash
# Delete local tags
git tag -d v1.2.3
git tag -d extra/fuegoecho/v1.2.3
git tag -d extra/fuegogin/v1.2.3
# ... etc

# Delete remote tags (if already pushed)
git push origin --delete v1.2.3
git push origin --delete extra/fuegoecho/v1.2.3
git push origin --delete extra/fuegogin/v1.2.3
# ... etc
```

### Getting Help

- **Local script**: `./scripts/release.sh --help`
- **GitHub Actions**: Check the workflow run logs for detailed information
- **Issues**: Report problems using the repository's issue tracker

## Examples

### Example 1: Standard Release via GitHub Actions

1. Go to Actions → Release → Run workflow
2. Enter `v1.3.0`
3. Leave "Dry run" unchecked
4. Click "Run workflow"
5. Wait for tests to pass
6. Approve the release when prompted
7. Check the created releases on GitHub

### Example 2: Testing with Local Script

```bash
# Preview what would be created
./scripts/release.sh --dry-run v1.3.0

# Create the release (with confirmation)
./scripts/release.sh v1.3.0

# Quick release (skip confirmation)
./scripts/release.sh --force v1.3.0
```

### Example 3: Pre-release

```bash
# Create a beta release
./scripts/release.sh v1.3.0-beta.1
```

## Migration from Manual Process

If you were previously creating tags manually:

1. **Stop** creating tags manually
2. **Use** either GitHub Actions or the local script
3. **Verify** that all expected tags are created
4. **Check** that GitHub releases are created properly

The automated system will handle all the complexity of multi-module releases for you!

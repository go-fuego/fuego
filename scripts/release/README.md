# Fuego Multi-Module Release System

Automated release tooling for managing synchronized releases across all 9 Fuego modules.

## Quick Start

```bash
# Preview a release (dry-run mode)
make release-dry-run

# Perform actual release
make release

# Show current versions
make release-versions

# Check if ready for release
make release-check
```

## Overview

This release system manages all 9 releasable modules in the Fuego project:
- Main module (`github.com/go-fuego/fuego`)
- CLI (`cmd/fuego`)
- 5 extra modules (`extra/fuegoecho`, `extra/fuegogin`, `extra/markdown`, `extra/sql`, `extra/sqlite3`)
- 2 middleware modules (`middleware/basicauth`, `middleware/cache`)

**All modules are tagged with the same version** for simplified versioning (e.g., v0.19.0).

## Make Targets

| Command | Description |
|---------|-------------|
| `make release` | Interactive release workflow with full validation |
| `make release-dry-run` | Preview mode - shows what would be tagged without making changes |
| `make release-versions` | Display current versions for all modules |
| `make release-validate` | Run full validation (lint + test + build) |
| `make release-check` | Quick check if ready for release (environment + git state) |
| `make release-rollback` | Emergency rollback - delete locally created tags (before push) |

## Release Workflow

### 1. Full Release (Interactive)

```bash
$ make release
```

This will:
1. **Validate environment** - Check Git, Go, golangci-lint are installed
2. **Check git state** - Ensure working directory is clean, on main branch, up-to-date
3. **Run comprehensive validation** - Execute lint, tests, and builds for all modules
4. **Interactive version selection** - Choose patch/minor/major or custom version
5. **Preview tags** - Show all 9 tags that will be created
6. **Create tags** - Create all tags locally (with confirmation)
7. **Push tags** - Push all tags to GitHub (with confirmation)

### 2. Dry-Run Mode

```bash
$ make release-dry-run
```

Perfect for:
- Testing the release process
- Previewing what tags would be created
- Validating version numbers
- CI/CD integration

### 3. Automated/CI Mode

```bash
# Skip all confirmations
export FUEGO_RELEASE_AUTO_CONFIRM=1

# Use specific version
./scripts/release/release.sh --version v0.19.0

# Skip validation (not recommended)
./scripts/release/release.sh --skip-validation
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `FUEGO_RELEASE_AUTO_CONFIRM=1` | Skip all confirmation prompts | `0` |
| `FUEGO_RELEASE_SKIP_LINT=1` | Skip linting | `0` |
| `FUEGO_RELEASE_SKIP_TESTS=1` | Skip tests | `0` |
| `FUEGO_RELEASE_SKIP_BUILD=1` | Skip builds | `0` |
| `FUEGO_RELEASE_ALLOW_NON_MAIN=1` | Allow release from non-main branch | `0` |
| `DEBUG=1` | Enable debug logging | `0` |

## Safety Features

✅ **Multiple validation layers** - Environment, git state, lint, tests, builds
✅ **Atomic tag creation** - All 9 tags created or none (rollback on failure)
✅ **Confirmation prompts** - Explicit confirmation before creating and pushing tags
✅ **Dry-run mode** - Preview without making changes
✅ **Emergency rollback** - Delete local tags if something goes wrong
✅ **Retracted version protection** - Blocks v1.0.0 and v1.0.1 (see go.mod)
✅ **Duplicate detection** - Prevents creating tags that already exist

## Versioning

### Version Format

- Semantic versioning: `vX.Y.Z` (e.g., v0.19.0)
- All 9 modules use the **same version**
- Patch: Bug fixes (v0.18.8 → v0.18.9)
- Minor: New features (v0.18.8 → v0.19.0)
- Major: Breaking changes (v0.18.8 → v1.0.0)

### Retracted Versions

The following versions are retracted and cannot be used:
- **v1.0.0** - Published accidentally
- **v1.0.1** - Contains retractions only

The release tool automatically skips these versions.

## Tagging Pattern

| Module | Tag Format | Example |
|--------|-----------|---------|
| Main | `vX.Y.Z` | `v0.19.0` |
| cmd/fuego | `cmd/fuego/vX.Y.Z` | `cmd/fuego/v0.19.0` |
| extra/* | `extra/{name}/vX.Y.Z` | `extra/fuegoecho/v0.19.0` |
| middleware/* | `middleware/{name}/vX.Y.Z` | `middleware/cache/v0.19.0` |

## Troubleshooting

### Uncommitted Changes

```
✗ Working directory has uncommitted changes
```

**Fix:** Commit or stash your changes:
```bash
git status
git add .
git commit -m "Your commit message"
```

### Not Up-to-Date

```
✗ Not up-to-date with origin/main
```

**Fix:** Pull latest changes:
```bash
git pull origin main
```

### Unpushed Commits

```
✗ You have unpushed commits
```

**Fix:** Push your commits:
```bash
git push origin main
```

### Test/Lint/Build Failures

```
✗ Tests failed
```

**Fix:** Run the checks manually to see detailed output:
```bash
make test    # Run tests
make lint    # Run linter
make build   # Run builds
```

### Tag Already Exists

```
✗ Version v0.19.0 already exists
```

**Fix:** Choose a different version or delete the existing tag (if you're sure):
```bash
# List existing tags
git tag -l "v0.19*"

# Delete tag locally and remotely (be careful!)
git tag -d v0.19.0
git push origin :refs/tags/v0.19.0
```

### Emergency Rollback

If tag creation fails mid-process:

```bash
make release-rollback
```

This deletes all locally created tags that haven't been pushed to remote.

## Scripts

The release system consists of several scripts:

- **`release.sh`** - Main orchestrator
- **`version.sh`** - Version detection and management
- **`validate.sh`** - Pre-release validation
- **`tag.sh`** - Tag creation and pushing
- **`lib/common.sh`** - Shared utilities (logging, colors, helpers)
- **`lib/config.sh`** - Module configuration

## Future Enhancements

Potential additions for Phase 2:

- GitHub releases creation (`gh release create`)
- Automated changelog generation
- Update go.mod dependencies across modules
- Release announcements (Discord, Twitter)
- CI/CD GitHub Actions workflow

## Examples

### Example 1: Standard Release

```bash
$ make release

🚀 Fuego Multi-Module Release
==============================

Running pre-release checks...
✓ All checks passed!

Select version:
  1) v0.18.9 (patch)
  2) v0.19.0 (minor) [recommended]
  3) v1.0.0 (major)

Choice [2]: 2

Release Summary
===============
Version: v0.19.0
Tags to create: 9 tags

Proceed? [y/N] y

Creating tags...
✓ All 9 tags created locally

Push to remote? [y/N] y
✓ Tags pushed to origin

🎉 Release v0.19.0 Complete!
```

### Example 2: Automated CI Release

```bash
#!/bin/bash
export FUEGO_RELEASE_AUTO_CONFIRM=1
./scripts/release/release.sh --version v0.19.0
```

### Example 3: Check Before Release

```bash
$ make release-check

Environment Checks
==================
✓ Git is installed
✓ Go is installed
✓ golangci-lint is installed
✓ In git repository
✓ At repository root
✓ All 9 modules found

Git State Checks
================
✓ Working directory is clean
✓ On main branch
✓ Up-to-date with origin/main
✓ No unpushed commits

✅ All pre-release checks passed!
```

## Support

For issues or questions:
- Check this README
- Review error messages (they include remediation steps)
- Check the scripts with `--help` flag
- File an issue on GitHub

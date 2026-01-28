# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Current Version**: v0.1.2
**Status**: Active Development
**Language**: Go 1.24.11

This is a **minimalist command-line bookmark management tool** written in Go. It provides a simple, opinionated interface for creating, organizing, and managing directory bookmarks using symbolic links, with multi-shell completion support (Bash, Zsh, Fish).

## Core Architecture

### Single Binary Design

- **Language**: Go (1.24.11+)
- **Main files**: `main.go` (core functionality), `completion.go` (shell completion)
- **Philosophy**: Unix philosophy - do one thing well
- **Dependencies**: Zero external dependencies, just Go standard library
- **Storage**: Symbolic links in `~/.marks/` directory

### Key Features

- **Symbolic link bookmarks** (stored in `~/.marks/` directory)
- **Simple bookmark creation** (`mark` uses current dir name, `mark <name>` uses custom name)
- **Multi-shell tab completion** (Bash, Zsh, Fish with `--autocomplete` flag)
- **Simple setup** (only asks for marks directory location)
- **Jump functionality** (`mark -j <name>` prints path for shell function wrapper)
- **Broken symlink detection** (marks broken bookmarks in list output)
- **Shell aliases** (`marks`, `unmark`, `jump` via `--alias` command)
- **Zero dependencies** (single static binary, no external libraries)
- **Configuration** stored in `~/.mark` file
- **Version tracking** built into release binaries (`--version` flag)

### Differences from Note Template

Mark is based on the `note` project template but simplified for bookmark management:

**Removed features:**
- No editor configuration (mark doesn't edit files)
- No archive system (bookmarks are just deleted)
- No search functionality (bookmarks are just names)
- No date stamping (bookmarks use simple names)
- No markdown file handling

**New features:**
- Symlink creation and management
- Jump functionality for directory navigation
- Broken symlink detection

## Development Commands

### Available Make Targets

```bash
# Build Commands
make build          # Build the binary
make release        # Build release binary with version info
make install        # Install system-wide (requires sudo)
make clean          # Clean build artifacts

# Testing Commands
make test                 # Run all test suites (unit, integration, completion, setup)
make test-unit            # Run unit tests only
make integration-test     # Run integration tests
make completion-test      # Run tab completion tests
make setup-test           # Run setup/configuration tests

# Code Quality Commands
make fmt            # Format Go code
make vet            # Run Go static analysis

# Release Commands
make bump           # Bump patch version
make bump-minor     # Bump minor version
make bump-major     # Bump major version
```

### Test Commands

```bash
# Unit tests
make test-unit

# Integration tests
make integration-test

# Completion tests
make completion-test

# Setup tests
make setup-test

# All tests
make test
```

## Validation Commands

When working on this project, **ALWAYS** run these validation steps before committing:

```bash
# Level 0: Clean and Check (CRITICAL)
make clean                  # Remove build artifacts
make vet                    # Static analysis
make fmt                    # Format code
git diff --exit-code        # Verify fmt made no changes

# Level 1: Build Check
make build                  # Build the binary
./mark --version            # Verify build

# Level 2: Unit Tests
make test-unit

# Level 3: All Tests (REQUIRED before release)
make test
```

## Project Structure

```
mark/
├── main.go                       # Main application code (bookmark management)
├── completion.go                 # Shell completion (bash/zsh/fish)
├── main_test.go                  # Unit tests
├── go.mod                        # Go module definition
├── Makefile                      # Build automation and release management
├── README.md                     # User documentation
├── RELEASE.md                    # Release notes and version history
├── CLAUDE.md                     # This file - guidance for Claude Code
├── COPYING.md                    # GPL-3.0 license
├── .gitignore                    # Ignore build artifacts
└── scripts/                      # Test and utility scripts
    ├── integration_test.sh           # Integration tests
    ├── completion_test.sh            # Tab completion tests
    └── setup_integration_test.sh     # Setup tests
```

## Development Guidelines

### CRITICAL REQUIREMENT: Never Run `git add`

**NEVER** run `git add`, `git add -A`, `git add .`, or `git add <file>` unless:
1. Adding release notes (RELEASE.md) during version bumping
2. Explicitly instructed by the user to stage specific files

The user stages files manually as a review checkpoint before commits. Running `git add` bypasses their review process. Always wait for the user to stage changes themselves.

**Exception:** During automated release workflows (`/development:release`), staging is permitted as part of the version bump process.

### CRITICAL REQUIREMENT: Test-Driven Development

**MANDATORY**: When adding functionality, modifying functionality, or fixing bugs, you **MUST ALWAYS**:

1. **Write tests** to verify the new functionality works correctly
2. **Add regression tests** for bug fixes to prevent the bug from reoccurring
3. **Update existing tests** when modifying functionality
4. **Run `make test`** to ensure all tests pass before committing

This is **NON-NEGOTIABLE**. No code changes should be committed without corresponding test coverage.

Test placement:
- **Unit tests** → `main_test.go` for core functions
- **Integration tests** → `scripts/integration_test.sh` for end-to-end workflows
- **Completion tests** → `scripts/completion_test.sh` for shell completion
- **Setup tests** → `scripts/setup_integration_test.sh` for configuration

### Code Patterns

- Single-file architecture in `main.go` with `completion.go` for shell integration
- Struct-based configuration (`Config` type)
- Symlink-based bookmark storage
- Comprehensive error handling
- File operations use `filepath` package for cross-platform compatibility

### Testing Strategy

The project includes automated tests across four test suites (~65 tests total):

| Suite | File | Tests | Coverage |
|-------|------|-------|----------|
| Unit | `main_test.go` | ~14 | Core functions, flag parsing, path handling, RC generation |
| Integration | `scripts/integration_test.sh` | ~13 | End-to-end workflows, bookmarking, jumping, broken links |
| Completion | `scripts/completion_test.sh` | ~29 | Tab completion for Bash/Zsh/Fish, partial matching, aliases |
| Setup | `scripts/setup_integration_test.sh` | ~11 | First-run setup, shell config, unified RC files |

**Unit Tests** (`main_test.go`):
- Core functionality, path handling, configuration
- Symlink creation and deletion
- Flag parsing
- RC file generation (bash, zsh, fish)
- Source line detection and feature parsing

**Integration Tests** (`scripts/integration_test.sh`):
- End-to-end user workflows
- Bookmark creation, listing, deletion, jumping
- Broken symlink handling
- Custom path bookmarks

**Completion Tests** (`scripts/completion_test.sh`):
- Tab completion for Bash, Zsh, Fish
- Partial matching
- Flag completion
- Alias completions (marks, unmark, jump)
- Broken bookmark formatting

**Setup Tests** (`scripts/setup_integration_test.sh`):
- First-run setup flow
- Configuration management
- Shell detection and alias setup
- Unified RC file creation
- Source line installation

### Key Functions to Understand

- `createBookmark()` - Creates symbolic links in ~/.marks/
- `listBookmarks()` - Lists bookmarks and shows broken symlinks
- `deleteBookmark()` - Removes bookmark symlinks
- `jumpBookmark()` - Resolves symlink and prints target path
- `expandPath()` - Handles tilde expansion and symlink resolution
- `parseFlags()` - Command-line flag parsing
- `setupAliases()` - Shell alias installation for marks, unmark, jump
- `RunAutocompleteSetup()` - Multi-shell completion setup

## Release Process

### Automated Release Workflow

Use the `/development:release` command for fully automated releases (if .claude/ commands are set up):

```bash
# Patch version bump (0.0.1 -> 0.0.2)
/development:release
/development:release patch

# Minor version bump (0.0.1 -> 0.1.0)
/development:release minor

# Major version bump (0.0.1 -> 1.0.0)
/development:release major
```

### Manual Release

```bash
# 1. Clean and validate
make clean && make vet && make fmt
git diff --exit-code  # Verify no fmt changes
make test         # All tests must pass

# 2. Bump version and tag
make bump             # Creates tag (e.g., v0.0.2)

# 3. Build release
make release          # Builds with version info

# 4. Validate binary
./mark --version      # Verify version, date, SHA

# 5. Push tag
git push origin v0.0.2
```

### Version Numbering

- **Patch** (0.0.X): Bug fixes, minor improvements, documentation updates
- **Minor** (0.X.0): New features, significant improvements (backward compatible)
- **Major** (X.0.0): Breaking changes, major rewrites

### Release Checklist

- [ ] All tests passing
- [ ] Code formatted (`make fmt`)
- [ ] No static analysis warnings (`make vet`)
- [ ] RELEASE.md updated with new version's release notes
- [ ] Binary validated (version, date, SHA correct)
- [ ] Tag pushed to GitHub

## Development Philosophy

Remember: This is a focused CLI tool following Unix philosophy. Keep changes minimal, well-tested, and true to the simple bookmark management workflow.

### Principles

- **Simplicity over features**: Only add what's truly needed
- **Test everything**: All changes require tests
- **No external dependencies**: Keep using only Go standard library
- **Cross-platform**: Test on Linux, macOS, BSD
- **Backward compatibility**: Don't break existing workflows
- **Documentation first**: Update docs before committing

### When Making Changes

1. **Read the code first**: Understand existing patterns
2. **Write tests first**: Test-driven development
3. **Keep it simple**: Resist feature creep
4. **Run all tests**: Use `make test` before committing
5. **Update docs**: README.md, RELEASE.md, and this file
6. **Follow conventions**: Match existing code style

## Core Operations

### Creating Bookmarks

```bash
mark                     # Creates bookmark with current directory name
mark myproject           # Creates bookmark named "myproject" pointing to current dir
mark work ~/projects     # Creates bookmark "work" pointing to ~/projects
```

Internally:
- If custom path provided, validates and expands the path (supports tilde expansion)
- Otherwise gets current directory with `os.Getwd()`
- Sanitizes name (replaces spaces with underscores)
- Creates symlink in `~/.marks/`
- Checks for duplicates

### Listing Bookmarks

```bash
mark -l
```

Output format:
```
downloads -> /home/user/Downloads
mark      -> /home/user/Projects/mark
oldproject -> [broken] /home/user/deleted-project
```

Internally:
- Reads directory entries from `~/.marks/`
- Checks each for symlink status
- Resolves targets and detects broken links
- Sorts alphabetically

### Deleting Bookmarks

```bash
mark -d myproject
```

Internally:
- Verifies bookmark exists
- Verifies it's a symlink (safety check)
- Removes symlink with `os.Remove()`

### Jumping to Bookmarks

```bash
mark -j myproject   # Prints: /home/user/Projects/myproject
jump myproject      # Shell function wraps mark -j and does cd
```

Internally (mark -j):
- Resolves symlink to target path
- Verifies target exists and is a directory
- Prints absolute path to stdout

Shell function (created by --alias):
```bash
jump() {
    local target=$(mark -j "$@")
    if [ $? -eq 0 ] && [ -n "$target" ]; then
        cd "$target"
    fi
}
```

## Common Tasks

### Adding a new flag

1. Add field to `ParsedFlags` struct
2. Update `parseFlags()` function to handle the flag
3. Add handling logic in `main()` function
4. Update `printHelp()` text
5. Add tests for the new flag
6. Update README.md with examples

### Modifying bookmark operations

1. Update the relevant function (`createBookmark`, `listBookmarks`, etc.)
2. Add unit tests in `main_test.go`
3. Add integration tests in `scripts/integration_test.sh`
4. Test edge cases (permissions, broken symlinks, etc.)
5. Update documentation

### Adding shell completion

1. Modify completion scripts in `completion.go`
2. Test with bash, zsh, and fish
3. Add tests to `scripts/completion_test.sh`
4. Document in README.md

## Troubleshooting

### Build Issues

- Check Go version: `go version` (need 1.24.11+)
- Clean and rebuild: `make clean && make build`
- Check for syntax errors: `make vet`

### Test Failures

- Run individual test suites to isolate issues
- Check for leftover test artifacts in `/tmp`
- Verify shell configuration files aren't corrupted

### Symlink Issues

- Check directory permissions: `ls -la ~/.marks`
- Verify symlinks with: `ls -la ~/.marks/<bookmark-name>`
- Test symlink creation manually: `ln -s /path/to/target ~/.marks/test`

### Shell Completion Issues

If `mark --autocomplete` doesn't work, manual setup:

**Bash** - Add to `~/.bashrc`:
```bash
[ -f ~/.mark_bash_rc ] && source ~/.mark_bash_rc
```

**Zsh** - Add to `~/.zshrc`:
```bash
[ -f ~/.mark_zsh_rc ] && source ~/.mark_zsh_rc
```

**Fish** - File should auto-load from `~/.config/fish/conf.d/mark.fish`

To regenerate RC files manually, delete the existing file and re-run `mark --autocomplete` or `mark --alias`.

## Related Projects

- [note](https://github.com/brockers/note) - A minimalist command line note creation/management tool (mark is based on this template)

## Future Enhancements (Maybe)

- Categories/tags for bookmarks
- Import/export bookmarks
- Bookmark descriptions
- Recently used bookmarks
- Fuzzy bookmark name matching
- Bookmark search by target path

Remember: Only add features that truly enhance the core bookmark workflow. Resist feature creep!

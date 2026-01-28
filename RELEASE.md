# Release Notes

## v0.1.3 - 2026-01-28

### Refactoring

- **Unified shell configuration** (df84aad)
  - Consolidate all shell configuration (aliases AND completions) into single RC files
  - Bash: `~/.mark_bash_rc`, Zsh: `~/.mark_zsh_rc`, Fish: `~/.config/fish/conf.d/mark.fish`
  - Single source line in shell config with existence check
  - Feature tracking header in RC files (`# Features: aliases completions`)
  - Automatic migration from legacy configuration format

### Documentation

- **Simplified README** (6cbbe40)
  - Reduced from 269 to 79 lines (70% smaller)
  - Quick Start workflow replaces verbose examples
  - Compact usage table for command reference
  - Consolidated developer docs into CLAUDE.md

- **Developer documentation improvements** (63231ee, 88b8ae0, f955260)
  - Added critical requirement to never run `git add` without permission
  - Expanded test coverage breakdown table
  - Manual shell completion troubleshooting guide
  - Updated version references

### Bug Fixes

- **Release script** (3ad9da0)
  - Fix incorrect binary name (`./note` → `./mark`)
  - Update test command (`make test-all` → `make test`)

### Features

- **New flags** (df84aad)
  - Add `--configure` as alias for `--config`
  - Update completion scripts to include `--configure`

---

## v0.1.2 - 2026-01-13

### Features

- **Completion enhancements**
  - Display bookmark paths in tab completion for easier navigation (f3a1156)
  - Improved visibility of bookmark targets during completion

- **Bookmark improvements**
  - Add custom path support for creating bookmarks (6a48409)
  - Enhanced broken bookmark visibility in listings

### Bug Fixes

- **Completion fixes**
  - Fix tab completion hanging and improve formatting consistency (6aa04f2)
  - Fix autocomplete not working for shell aliases (3ce6339)

### Refactoring

- **Project structure**
  - Reorganize Claude command files into commands directory (e551daf)
  - Better organization following Claude Code conventions

---

## v0.0.1 - Initial Release (TBD)

**Status**: In Development

### Features

- **Core bookmark management**
  - Create bookmarks with `mark` (uses current directory name)
  - Create named bookmarks with `mark <name>`
  - List all bookmarks with `mark -l`
  - Delete bookmarks with `mark -d <name>`
  - Jump to bookmarks with `mark -j <name>` (prints path for shell function)

- **Symlink-based storage**
  - Bookmarks stored as symbolic links in `~/.marks/`
  - Broken symlink detection and indication in listing
  - Simple, filesystem-native approach

- **Shell integration**
  - Bash/Zsh/Fish completion support
  - Shell aliases: `marks`, `unmark`, `jump`
  - Jump function wraps `mark -j` with cd command

- **Configuration**
  - Simple setup asks only for marks directory location
  - Config stored in `~/.mark` file
  - Automatic directory creation

- **Unix-style CLI**
  - Standard flags: `-h`, `-v`, `-l`, `-d`, `-j`
  - Long flags: `--help`, `--version`, `--config`, `--autocomplete`, `--alias`
  - Version information in binary

### Implementation Details

- Zero external dependencies (Go standard library only)
- Single binary design (mark binary only)
- GPL-3.0 licensed
- Cross-platform support (Linux, macOS, BSD)

### Planned Improvements

- Unit tests and integration test suite
- Comprehensive completion tests
- Setup integration tests
- Enhanced error handling for edge cases

---

## Release Process

### Creating a New Release

1. **Pre-release validation**
   ```bash
   make clean
   make vet
   make fmt
   git diff --exit-code  # Ensure no changes from fmt
   make test-all         # All tests must pass
   ```

2. **Update version**
   ```bash
   # For patch release (0.0.1 -> 0.0.2)
   make bump

   # For minor release (0.0.1 -> 0.1.0)
   make bump-minor

   # For major release (0.0.1 -> 1.0.0)
   make bump-major
   ```

3. **Build release binary**
   ```bash
   make release
   ./mark --version  # Verify version info
   ```

4. **Update RELEASE.md**
   - Add release notes for the new version
   - Document new features, bug fixes, changes
   - Include breaking changes if any

5. **Push tag**
   ```bash
   git push origin v0.0.2
   ```

### Version Numbering

Following semantic versioning (semver):
- **MAJOR** version: Incompatible API changes
- **MINOR** version: Backwards-compatible new functionality
- **PATCH** version: Backwards-compatible bug fixes

### Release Checklist

- [ ] All tests passing (`make test-all`)
- [ ] Code formatted (`make fmt`)
- [ ] No static analysis warnings (`make vet`)
- [ ] RELEASE.md updated with new version notes
- [ ] Version tag created (`make bump*`)
- [ ] Release binary built (`make release`)
- [ ] Binary version verified (`./mark --version`)
- [ ] Tag pushed to origin

---

## Changelog Format

Each release should document:

### Added
- New features

### Changed
- Changes to existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security improvements or fixes

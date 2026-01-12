# Release Notes

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

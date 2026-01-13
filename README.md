# mark

[![Version](https://img.shields.io/badge/version-0.1.2-blue.svg)](https://github.com/brockers/mark/releases/tag/v0.1.2)
[![License](https://img.shields.io/badge/license-GPL--3.0-green.svg)](https://www.gnu.org/licenses/gpl-3.0.en.html)
[![Tests](https://img.shields.io/badge/tests-48%20passing-brightgreen.svg)](#development)

A minimalist command-line bookmark management tool written in Go. Create, organize, and jump to directory bookmarks with tab completion and zero lock-in - just symbolic links.

## Features

- **Symlink-Based Bookmarks**: All bookmarks stored as symbolic links in `~/.marks/` directory
- **Simple Creation**: `mark` uses current directory name, `mark <name>` creates custom bookmark
- **Custom Paths**: Create bookmarks pointing to any directory with `mark <name> <path>`
- **Tab Completion**: Type `mark -j proj` and press TAB to see matching bookmarks (Bash, Zsh, Fish)
- **Quick Jump**: Navigate to bookmarked directories with `mark -j <name>` or the `jump` alias
- **Broken Link Detection**: Automatically detects and marks broken bookmarks in red
- **Shell Aliases**: Optional convenient shortcuts (`marks`, `unmark`, `jump`)
- **Two-Question Setup**: Simple configuration when first running the app
- **Zero Dependencies**: Single static binary, no external libraries required
- **Well Documented**: Standard Unix help via `-h` or `--help` options
- **Thoroughly Tested**: 48 automated tests covering unit, integration, completion, and setup scenarios

## Examples

### Create Default

Create a new bookmark in of the current folders

```bash
mark 
```

This will create a symbolic link in your ~/.marks/ folder, named the same as the folder you are in, that points to the folder you are in.

### Create Named

Create a new bookmark but with a specified name

```bash
mark downloads
```

Creates a symbolic link in your ~/.marks/ folder, named downloads that points to the folder you are currently in.

### Create with Custom Path

Create a bookmark pointing to a custom directory (not current directory)

```bash
mark work ~/projects/work
mark tmp /tmp
```

Creates a symbolic link in your ~/.marks/ folder with the specified name that points to the given path. The path can be absolute or relative, and tilde expansion is supported.

### Show Bookmarks

List all of your bookmarks and where they point to:

```bash
mark -l

  downloads -> /home/jsmith/Downloads
  mark      -> /home/jsmith/Project/mark
```

Cleanly displays all of the symbolic links in your ~/.marks/ folder.

### Delete Bookmark

```bash 
mark -d downloads 
```

Removes the sybmolic link in your ~/.marks/ folder named downloads.

### Go to Bookmark 

Jump to your bookmarked folder

```bash
mark -j downloads
```

Does a `cd ~/.marks/downloads` to send you to the named bookmark.

### Alias

mark also has a couple built in aliases including

```bash
marks  #same as mark -l 
unmark #same as mark -d 
jump   #same as mark -j
```

### Autocomplete

Finally mark has built in autocomplete so you can alway double tab to see which mark you will jump to or delete.

## Installation

### From Release Binary

Download the latest release from [GitHub Releases](https://github.com/brockers/mark/releases):

```bash
# Download the release binary (replace with latest version)
wget https://github.com/brockers/mark/releases/download/v0.1.2/mark

# Make it executable
chmod +x mark

# Move to your PATH
sudo mv mark /usr/local/bin/

# Verify installation
mark --version
```

### From Source

Requirements:
- Go 1.24 or later

```bash
# Clone the repository
git clone https://github.com/brockers/mark.git
cd mark

# Build the binary
make build

# Install system-wide (optional, requires sudo)
make install

# Or copy manually to your PATH
cp mark ~/bin/  # or wherever you keep personal binaries
```

### Enable Tab Completion

After installation, enable tab completion for your shell:

```bash
# Automatic setup (recommended)
mark --autocomplete

# This will configure completion for your detected shell (Bash, Zsh, or Fish)
```

Manual setup if needed:

```bash
# Bash
echo 'source <(mark --autocomplete bash)' >> ~/.bashrc

# Zsh
echo 'source <(mark --autocomplete zsh)' >> ~/.zshrc

# Fish
mark --autocomplete fish > ~/.config/fish/completions/mark.fish
```

### Enable Shell Aliases

Optionally set up convenient aliases for common commands:

```bash
mark --alias
```

This creates:
- `marks` - Same as `mark -l`
- `unmark` - Same as `mark -d`
- `jump` - Same as `mark -j` (includes cd wrapper function)

## Development

### Building

```bash
# Build the binary
make build

# Build release version with version info
make release

# Format code
make fmt

# Run static analysis
make vet

# Clean build artifacts
make clean
```

### Testing

`mark` has a comprehensive test suite with 48 automated tests:

```bash
# Run unit tests
make test

# Run integration tests
make integration-test

# Run completion tests
make completion-test

# Run setup tests
make setup-test

# Run all tests
make test-all
```

Test coverage includes:
- **Unit Tests** (5 tests): Core functionality, path handling, configuration
- **Integration Tests** (13 tests): End-to-end workflows, bookmarking, jumping, deletion
- **Completion Tests** (29 tests): Tab completion for Bash, Zsh, Fish, including aliases
- **Setup Tests** (6 tests): First-run setup, configuration, shell integration

### Release Process

To create a new release:

```bash
# Run the automated release workflow
# /development:release [patch|minor|major]

# Or manually:
make bump && make release && git push origin <TAG>
```

## Philosophy

`mark` follows the Unix philosophy: do one thing well and compose with other tools. It's intentionally minimal and opinionated to provide a frictionless bookmark experience for terminal users.

- **No databases**: Just symbolic links in a folder
- **No sync built-in**: Use git, Dropbox, or any sync tool you prefer
- **No categories or tags**: Use your filesystem organization
- **No dependencies**: Single static binary
- **No lock-in**: Your bookmarks are just symlinks

## Support & Contributing

- **Repository**: https://github.com/brockers/mark
- **Issues**: https://github.com/brockers/mark/issues
- **Releases**: https://github.com/brockers/mark/releases
- **Documentation**: See [RELEASE.md](RELEASE.md) for release notes and [CLAUDE.md](CLAUDE.md) for development guidance

Contributions are welcome! Please feel free to submit issues or pull requests.

## Related Projects

- [note](https://github.com/brockers/note) - A minimalist, opinionated command line note creation/management application.

## License

This program is free software licensed under GPL-3.0.
See https://www.gnu.org/licenses/ for details.

---

**Version 0.1.2** | Built with Go | [View Release Notes](RELEASE.md)


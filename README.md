# mark

[![License](https://img.shields.io/badge/license-GPL--3.0-green.svg)](https://www.gnu.org/licenses/gpl-3.0.en.html)

A minimalist command-line bookmark manager. Create, list, and jump to directory bookmarks using symbolic links.

## Quick Start

```bash
# Bookmark your current directory
cd ~/projects/myapp
mark                    # creates bookmark named "myapp"

# Or with a custom name
mark work               # creates bookmark named "work"

# List all bookmarks
mark -l
#   myapp -> /home/user/projects/myapp
#   work  -> /home/user/projects/myapp

# Jump to a bookmark
jump work               # cd's to the bookmarked directory

# Delete a bookmark
mark -d work
```

On first run, `mark` will prompt to set up tab completion and shell aliases (`marks`, `unmark`, `jump`).

## Installation

**From source:**
```bash
git clone https://github.com/brockers/mark.git
cd mark
make build
sudo cp mark /usr/local/bin/   # or: cp mark ~/bin/
```

**From release:** Download from [GitHub Releases](https://github.com/brockers/mark/releases)

## Usage

| Command | Description |
|---------|-------------|
| `mark` | Bookmark current directory using folder name |
| `mark <name>` | Bookmark current directory with custom name |
| `mark <name> <path>` | Bookmark a specific path |
| `mark -l` | List all bookmarks |
| `mark -d <name>` | Delete a bookmark |
| `mark -j <name>` | Print bookmark path (used by `jump`) |
| `mark --config` | Re-run setup (completion, aliases) |

**Aliases** (after running `mark --alias`):
- `marks` → `mark -l`
- `unmark` → `mark -d`
- `jump` → `mark -j` with `cd`

## Philosophy

- **No databases** — just symlinks in `~/.marks/`
- **No lock-in** — bookmarks are plain symbolic links
- **No dependencies** — single static binary
- **No sync** — use git, Dropbox, or any tool you prefer

## Development

```bash
make build      # build binary
make test       # run all tests
make help       # see all targets
```

## License

GPL-3.0 — See [COPYING.md](COPYING.md)

**Links:** [Repository](https://github.com/brockers/mark) · [Issues](https://github.com/brockers/mark/issues) · [Releases](https://github.com/brockers/mark/releases)

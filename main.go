/*
Copyright (C) 2025  Mark CLI Contributors

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type Config struct {
	MarksDir string
}

var (
	Version   = "dev"
	CommitSHA = "not set"
	BuildDate = "not set"
)

const (
	// ANSI color codes
	colorRed   = "\033[0;31m"
	colorReset = "\033[0m"
)

func main() {
	// Parse custom flags with Unix-like behavior first
	flags, args := parseFlags(os.Args[1:])

	// Handle version number (before config load)
	if flags.Version {
		printVersion()
		return
	}

	// Handle help (before config load)
	if flags.Help {
		printHelp()
		return
	}

	// Load config after checking version/help
	config, firstTimeSetup := loadOrCreateConfig()

	// If first-time setup was just completed, exit gracefully
	if firstTimeSetup {
		return
	}

	// Handle config
	if flags.Config {
		runSetup()
		os.Exit(0)
	}

	// Handle autocomplete setup
	if flags.Autocomplete {
		RunAutocompleteSetup()
		return
	}

	// Handle alias setup
	if flags.Alias {
		RunAliasSetup()
		return
	}

	// Handle listing
	if flags.List {
		listBookmarks(config)
		return
	}

	// Handle delete
	if flags.Delete != "" {
		deleteBookmark(config, flags.Delete)
		return
	}

	// Handle jump
	if flags.Jump != "" {
		jumpBookmark(config, flags.Jump)
		return
	}

	// Handle bookmark creation
	bookmarkName := ""
	targetPath := ""

	if len(args) == 1 {
		// Single argument: bookmark name, use current directory as target
		bookmarkName = args[0]
	} else if len(args) >= 2 {
		// Two arguments: bookmark name and custom path
		bookmarkName = args[0]
		targetPath = args[1]
	}
	// else: no arguments, createBookmark will use current directory name

	createBookmark(config, bookmarkName, targetPath)
}

func loadOrCreateConfig() (Config, bool) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(homeDir, ".mark")

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// First run, create config
		return runSetup(), true
	}

	// Load existing config
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening config: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	config := Config{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "marksdir":
			config.MarksDir = expandPath(value)
		}
	}

	if config.MarksDir == "" {
		fmt.Println("Invalid config file. Running setup...")
		return runSetup(), false
	}

	return config, false
}

func runSetup() Config {
	reader := bufio.NewReader(os.Stdin)
	config := Config{}

	// Get current values if they exist
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".mark")
	if file, err := os.Open(configPath); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				switch strings.TrimSpace(parts[0]) {
				case "marksdir":
					config.MarksDir = expandPath(strings.TrimSpace(parts[1]))
				}
			}
		}
		file.Close()
	}

	// Ask for marks directory
	defaultDir := config.MarksDir
	if defaultDir == "" {
		defaultDir = "~/.marks"
	}

	fmt.Printf("Where should bookmarks be stored (%s): ", defaultDir)
	marksDir, _ := reader.ReadString('\n')
	marksDir = strings.TrimSpace(marksDir)
	if marksDir == "" {
		marksDir = defaultDir
	}

	marksDir = expandPath(marksDir)
	fmt.Printf("Setting your bookmarks location to %s ...\n", marksDir)
	config.MarksDir = marksDir

	// Create directory if it doesn't exist
	if err := os.MkdirAll(config.MarksDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating marks directory: %v\n", err)
		os.Exit(1)
	}

	// Ask about command line completion
	SetupCompletion(reader)

	// Ask about shell aliases
	setupAliases(reader)

	// Save config
	saveConfig(config)
	return config
}

func saveConfig(config Config) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(homeDir, ".mark")
	file, err := os.Create(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Convert absolute path back to ~ notation for config file
	marksDir := config.MarksDir
	if strings.HasPrefix(marksDir, homeDir) {
		marksDir = "~" + strings.TrimPrefix(marksDir, homeDir)
	}

	fmt.Fprintf(file, "marksdir=%s\n", marksDir)
}

func setupAliases(reader *bufio.Reader) {
	// Check if aliases are already set up
	if areAliasesAlreadySetup() {
		return
	}

	fmt.Println()
	fmt.Print("Would you like to set up shell aliases (marks, unmark, jump)? (y/N): ")
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "y" && response != "yes" {
		fmt.Println("Skipping alias setup. You can run 'mark --config' later to set them up.")
		return
	}

	shell := detectShell()
	if shell == "" {
		fmt.Println("Could not detect shell type. Skipping alias setup.")
		return
	}

	switch shell {
	case "bash":
		setupBashAliases()
	case "zsh":
		setupZshAliases()
	case "fish":
		setupFishAliases()
	default:
		fmt.Printf("Shell '%s' not supported for aliases. Supported shells: bash, zsh, fish\n", shell)
	}
}

func areAliasesAlreadySetup() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	shell := detectShell()
	switch shell {
	case "bash":
		bashrc := filepath.Join(homeDir, ".bashrc")
		if content, err := os.ReadFile(bashrc); err == nil {
			contentStr := string(content)
			return strings.Contains(contentStr, "alias marks=") && strings.Contains(contentStr, "alias unmark=") && strings.Contains(contentStr, "function jump")
		}
	case "zsh":
		zshrc := filepath.Join(homeDir, ".zshrc")
		if content, err := os.ReadFile(zshrc); err == nil {
			contentStr := string(content)
			return strings.Contains(contentStr, "alias marks=") && strings.Contains(contentStr, "alias unmark=") && strings.Contains(contentStr, "function jump")
		}
	case "fish":
		fishConfigDir := filepath.Join(homeDir, ".config", "fish", "config.fish")
		if content, err := os.ReadFile(fishConfigDir); err == nil {
			contentStr := string(content)
			return strings.Contains(contentStr, "alias marks ") && strings.Contains(contentStr, "alias unmark ") && strings.Contains(contentStr, "function jump")
		}
	}
	return false
}

func setupBashAliases() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		return
	}

	bashrcPath := filepath.Join(homeDir, ".bashrc")

	// Get the full path to the mark binary
	markPath, err := os.Executable()
	if err != nil {
		// Fallback to checking PATH
		markPath, err = exec.LookPath("mark")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not determine mark command path: %v\n", err)
			return
		}
	}

	aliasLines := fmt.Sprintf(`
# mark command aliases
alias marks='%s -l'
alias unmark='%s -d'
function jump() {
    local target=$(%s -j "$@")
    if [ $? -eq 0 ] && [ -n "$target" ]; then
        cd "$target"
    fi
}
`, markPath, markPath, markPath)

	// Append to .bashrc
	file, err := os.OpenFile(bashrcPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening .bashrc: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(aliasLines); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing aliases to .bashrc: %v\n", err)
		return
	}

	fmt.Printf("✓ Bash aliases setup complete!\n")
	fmt.Printf("  Added 'marks', 'unmark', and 'jump' aliases to %s\n", bashrcPath)
	fmt.Printf("  Run 'source ~/.bashrc' or restart your shell to activate aliases\n")
}

func setupZshAliases() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		return
	}

	zshrcPath := filepath.Join(homeDir, ".zshrc")

	// Get the full path to the mark binary
	markPath, err := os.Executable()
	if err != nil {
		// Fallback to checking PATH
		markPath, err = exec.LookPath("mark")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not determine mark command path: %v\n", err)
			return
		}
	}

	aliasLines := fmt.Sprintf(`
# mark command aliases
alias marks='%s -l'
alias unmark='%s -d'
function jump() {
    local target=$(%s -j "$@")
    if [ $? -eq 0 ] && [ -n "$target" ]; then
        cd "$target"
    fi
}
`, markPath, markPath, markPath)

	// Append to .zshrc
	file, err := os.OpenFile(zshrcPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening .zshrc: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(aliasLines); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing aliases to .zshrc: %v\n", err)
		return
	}

	fmt.Printf("✓ Zsh aliases setup complete!\n")
	fmt.Printf("  Added 'marks', 'unmark', and 'jump' aliases to %s\n", zshrcPath)
	fmt.Printf("  Run 'source ~/.zshrc' or restart your shell to activate aliases\n")
}

func setupFishAliases() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		return
	}

	// Create fish config directory if it doesn't exist
	fishConfigDir := filepath.Join(homeDir, ".config", "fish")
	if err := os.MkdirAll(fishConfigDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating fish config directory: %v\n", err)
		return
	}

	fishConfigPath := filepath.Join(fishConfigDir, "config.fish")

	// Get the full path to the mark binary
	markPath, err := os.Executable()
	if err != nil {
		// Fallback to checking PATH
		markPath, err = exec.LookPath("mark")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not determine mark command path: %v\n", err)
			return
		}
	}

	aliasLines := fmt.Sprintf(`
# mark command aliases
alias marks '%s -l'
alias unmark '%s -d'
function jump
    set -l target (%s -j $argv)
    if test $status -eq 0 -a -n "$target"
        cd "$target"
    end
end
`, markPath, markPath, markPath)

	// Append to config.fish
	file, err := os.OpenFile(fishConfigPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening fish config: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(aliasLines); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing aliases to fish config: %v\n", err)
		return
	}

	fmt.Printf("✓ Fish aliases setup complete!\n")
	fmt.Printf("  Added 'marks', 'unmark', and 'jump' aliases to %s\n", fishConfigPath)
	fmt.Printf("  Restart your shell to activate aliases\n")
}

func expandPath(path string) string {
	// Handle tilde expansion first
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		path = filepath.Join(homeDir, path[2:])
	}

	// Resolve symbolic links to get the actual path
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		// If we can't resolve symlinks, return the original path
		// This handles cases where the path doesn't exist yet or other errors
		return path
	}

	return resolvedPath
}

func createBookmark(config Config, name string, targetPath string) {
	var targetDir string

	// Determine target directory
	if targetPath != "" {
		// Custom path provided - expand and validate it
		targetDir = expandPath(targetPath)

		// Verify the target directory exists
		fileInfo, err := os.Stat(targetDir)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: Target directory does not exist: %s\n", targetPath)
			} else {
				fmt.Fprintf(os.Stderr, "Error accessing target directory: %v\n", err)
			}
			os.Exit(1)
		}

		// Verify it's a directory
		if !fileInfo.IsDir() {
			fmt.Fprintf(os.Stderr, "Error: Target path is not a directory: %s\n", targetPath)
			os.Exit(1)
		}
	} else {
		// No custom path - use current working directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		targetDir = currentDir
	}

	// If name is empty, use the target directory name
	if name == "" {
		name = filepath.Base(targetDir)
	}

	// Sanitize bookmark name
	// Replace spaces with underscores and remove path separators
	name = strings.ReplaceAll(name, " ", "_")
	if strings.Contains(name, string(os.PathSeparator)) {
		fmt.Fprintf(os.Stderr, "Error: Bookmark name cannot contain path separators\n")
		os.Exit(1)
	}

	if name == "" {
		fmt.Fprintf(os.Stderr, "Error: Bookmark name cannot be empty\n")
		os.Exit(1)
	}

	// Ensure marks directory exists
	if err := os.MkdirAll(config.MarksDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating marks directory: %v\n", err)
		os.Exit(1)
	}

	// Check if bookmark already exists
	symlinkPath := filepath.Join(config.MarksDir, name)
	if _, err := os.Lstat(symlinkPath); err == nil {
		fmt.Fprintf(os.Stderr, "Error: Bookmark '%s' already exists. Use 'mark -d %s' to remove it first.\n", name, name)
		os.Exit(1)
	}

	// Create the symlink
	if err := os.Symlink(targetDir, symlinkPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating bookmark: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Created bookmark '%s' -> %s\n", name, targetDir)
}

func listBookmarks(config Config) {
	// Read directory entries
	entries, err := os.ReadDir(config.MarksDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No bookmarks found. Create one with 'mark <name>'")
			return
		}
		fmt.Fprintf(os.Stderr, "Error reading bookmarks directory: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("No bookmarks found. Create one with 'mark <name>'")
		return
	}

	// Collect bookmark information
	type bookmarkInfo struct {
		name   string
		target string
		broken bool
	}

	var bookmarks []bookmarkInfo

	for _, entry := range entries {
		symlinkPath := filepath.Join(config.MarksDir, entry.Name())

		// Check if it's a symlink
		fileInfo, err := os.Lstat(symlinkPath)
		if err != nil {
			continue
		}

		if fileInfo.Mode()&os.ModeSymlink == 0 {
			// Not a symlink, skip
			continue
		}

		// Read symlink target
		target, err := os.Readlink(symlinkPath)
		if err != nil {
			continue
		}

		// Check if target exists
		_, err = os.Stat(symlinkPath)
		broken := err != nil

		bookmarks = append(bookmarks, bookmarkInfo{
			name:   entry.Name(),
			target: target,
			broken: broken,
		})
	}

	// Sort alphabetically by name
	sort.Slice(bookmarks, func(i, j int) bool {
		return bookmarks[i].name < bookmarks[j].name
	})

	// Print bookmarks
	for _, bm := range bookmarks {
		if bm.broken {
			fmt.Printf("  %s -> [%sbroken%s] %s%s%s\n", bm.name, colorRed, colorReset, colorRed, bm.target, colorReset)
		} else {
			fmt.Printf("  %s -> %s\n", bm.name, bm.target)
		}
	}
}

func deleteBookmark(config Config, name string) {
	if name == "" {
		fmt.Fprintf(os.Stderr, "Error: Bookmark name required for -d flag\n")
		os.Exit(1)
	}

	symlinkPath := filepath.Join(config.MarksDir, name)

	// Check if bookmark exists
	fileInfo, err := os.Lstat(symlinkPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Bookmark '%s' does not exist\n", name)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error accessing bookmark: %v\n", err)
		os.Exit(1)
	}

	// Verify it's a symlink
	if fileInfo.Mode()&os.ModeSymlink == 0 {
		fmt.Fprintf(os.Stderr, "Error: '%s' is not a bookmark (not a symlink)\n", name)
		os.Exit(1)
	}

	// Remove the symlink
	if err := os.Remove(symlinkPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing bookmark: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Removed bookmark '%s'\n", name)
}

func jumpBookmark(config Config, name string) {
	if name == "" {
		fmt.Fprintf(os.Stderr, "Error: Bookmark name required for -j flag\n")
		os.Exit(1)
	}

	symlinkPath := filepath.Join(config.MarksDir, name)

	// Check if bookmark exists
	fileInfo, err := os.Lstat(symlinkPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Bookmark '%s' does not exist\n", name)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error accessing bookmark: %v\n", err)
		os.Exit(1)
	}

	// Verify it's a symlink
	if fileInfo.Mode()&os.ModeSymlink == 0 {
		fmt.Fprintf(os.Stderr, "Error: '%s' is not a bookmark (not a symlink)\n", name)
		os.Exit(1)
	}

	// Resolve the symlink to get the actual target
	targetPath, err := filepath.EvalSymlinks(symlinkPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Bookmark '%s' points to non-existent directory\n", name)
		os.Exit(1)
	}

	// Verify target is a directory
	targetInfo, err := os.Stat(targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Bookmark '%s' points to non-existent directory\n", name)
		os.Exit(1)
	}

	if !targetInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: Bookmark '%s' points to a file, not a directory\n", name)
		os.Exit(1)
	}

	// Print the target path to stdout (for shell function to capture)
	fmt.Println(targetPath)
}

// ParsedFlags represents parsed command line flags
type ParsedFlags struct {
	List         bool
	Delete       string
	Jump         string
	Config       bool
	Autocomplete bool
	Alias        bool
	Help         bool
	Version      bool
}

// parseFlags implements Unix-like flag parsing
func parseFlags(args []string) (*ParsedFlags, []string) {
	flags := &ParsedFlags{}
	var remainingArgs []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--help" {
			flags.Help = true
		} else if arg == "--version" {
			flags.Version = true
		} else if arg == "--config" {
			flags.Config = true
		} else if arg == "--autocomplete" {
			flags.Autocomplete = true
		} else if arg == "--alias" {
			flags.Alias = true
		} else if strings.HasPrefix(arg, "--") {
			// Unknown long flag, treat as regular argument
			remainingArgs = append(remainingArgs, arg)
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// Handle short flags
			flagChars := arg[1:] // Remove the '-' prefix

			for j, char := range flagChars {
				switch char {
				case 'v':
					flags.Version = true
				case 'h':
					flags.Help = true
				case 'l':
					flags.List = true
				case 'd':
					// -d requires an argument
					if j == len(flagChars)-1 {
						// -d is the last flag, next arg is the bookmark name
						if i+1 < len(args) {
							i++
							flags.Delete = args[i]
						} else {
							fmt.Fprintf(os.Stderr, "Error: -d flag requires a bookmark name\n")
							os.Exit(1)
						}
					} else {
						fmt.Fprintf(os.Stderr, "Error: -d flag must be the last in a flag chain\n")
						os.Exit(1)
					}
				case 'j':
					// -j requires an argument
					if j == len(flagChars)-1 {
						// -j is the last flag, next arg is the bookmark name
						if i+1 < len(args) {
							i++
							flags.Jump = args[i]
						} else {
							fmt.Fprintf(os.Stderr, "Error: -j flag requires a bookmark name\n")
							os.Exit(1)
						}
					} else {
						fmt.Fprintf(os.Stderr, "Error: -j flag must be the last in a flag chain\n")
						os.Exit(1)
					}
				default:
					fmt.Fprintf(os.Stderr, "Error: unknown flag -%c\n", char)
					os.Exit(1)
				}
			}
		} else {
			// Regular argument
			remainingArgs = append(remainingArgs, arg)
		}
	}

	return flags, remainingArgs
}

// RunAliasSetup handles the standalone alias setup flow
func RunAliasSetup() {
	fmt.Println("mark - Shell Alias Setup")
	fmt.Println()
	fmt.Println("This will set up convenient shell aliases:")
	fmt.Println("• marks -> mark -l")
	fmt.Println("• unmark -> mark -d")
	fmt.Println("• jump -> mark -j (with cd wrapper)")
	fmt.Println()

	// Check if aliases are already set up
	if areAliasesAlreadySetup() {
		fmt.Println("Aliases are already set up!")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	// Use the existing setupAliases function for the core logic
	setupAliases(reader)
}

func printVersion() {
	fmt.Println(Version)
}

func printHelp() {
	fmt.Println(`mark - A minimalist CLI bookmark tool

USAGE:
  mark                 Create bookmark with current directory name
  mark <name>          Create bookmark with custom name
  mark <name> <path>   Create bookmark pointing to custom path
  mark [OPTIONS]

OPTIONS:
  -l                   List all bookmarks
  -d <name>            Delete bookmark
  -j <name>            Jump to bookmark (prints path)
  -h                   Show this help message
  -v                   Print version number

  --help               Show this help message
  --config             Run setup/reconfigure
  --autocomplete       Setup/update command line autocompletion
  --alias              Setup/update shell aliases
  --version            Print version number

EXAMPLES:
  mark                 Create bookmark (if in ~/projects, creates 'projects')
  mark downloads       Create bookmark 'downloads' pointing to current dir
  mark work ~/work     Create bookmark 'work' pointing to ~/work
  mark tmp /tmp        Create bookmark 'tmp' pointing to /tmp
  mark -l              List all bookmarks with their targets
  mark -d downloads    Delete the 'downloads' bookmark
  mark -j projects     Print path to 'projects' bookmark
  jump projects        Change directory to 'projects' (requires alias setup)

ALIASES:
  After running 'mark --alias', you can use:
  marks                Same as 'mark -l'
  unmark <name>        Same as 'mark -d <name>'
  jump <name>          Change directory to bookmark

CONFIGURATION:
  Settings are stored in ~/.mark
  Bookmarks are stored in ~/.marks/ as symbolic links
  Use 'mark --config' to reconfigure

RELEASE:
     Version:    ` + Version + `
  Build Date:    ` + BuildDate + `
  Commit SHA:    ` + CommitSHA + `

LICENSE:
  This program is free software licensed under GPL-3.0.
  See <https://www.gnu.org/licenses/> for details.

For more information, see: https://github.com/brockers/mark`)
}

// detectShell detects the current shell from environment variables
func detectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return ""
	}

	// Extract shell name from path
	shellName := filepath.Base(shell)

	// Map common shell variants
	switch shellName {
	case "bash":
		return "bash"
	case "zsh":
		return "zsh"
	case "fish":
		return "fish"
	default:
		return shellName
	}
}

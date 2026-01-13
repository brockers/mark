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
	"path/filepath"
	"strings"
)

// SetupCompletion handles the interactive completion setup prompt
func SetupCompletion(reader *bufio.Reader) {
	// Check if completion is already set up
	if IsCompletionAlreadySetup() {
		return
	}

	fmt.Println()
	fmt.Print("Would you like to set up command line completion for mark? (y/N): ")
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "y" && response != "yes" {
		fmt.Println("Skipping completion setup. You can run 'mark --config' later to set it up.")
		return
	}

	shell := detectShell()
	if shell == "" {
		fmt.Println("Could not detect shell type. Skipping completion setup.")
		return
	}

	switch shell {
	case "bash":
		SetupBashCompletion()
	case "zsh":
		SetupZshCompletion()
	case "fish":
		SetupFishCompletion()
	default:
		fmt.Printf("Shell '%s' not supported for completion. Supported shells: bash, zsh, fish\n", shell)
	}
}

// IsCompletionAlreadySetup checks if command line completion is already configured
func IsCompletionAlreadySetup() bool {
	shell := detectShell()
	if shell == "" {
		return false
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	switch shell {
	case "bash":
		// Check if ~/.mark.bash exists and is sourced in shell config
		bashCompletionFile := filepath.Join(homeDir, ".mark.bash")
		if _, err := os.Stat(bashCompletionFile); err == nil {
			// Check .bashrc or .bash_profile for mark completion
			bashFiles := []string{".bashrc", ".bash_profile", ".profile"}
			for _, file := range bashFiles {
				if CheckFileForCompletionSource(filepath.Join(homeDir, file)) {
					return true
				}
			}
		}
	case "zsh":
		zshCompletionFile := filepath.Join(homeDir, ".mark.zsh")
		if _, err := os.Stat(zshCompletionFile); err == nil {
			if CheckFileForCompletionSource(filepath.Join(homeDir, ".zshrc")) {
				return true
			}
		}
	case "fish":
		// Check fish completion directory
		fishCompletionDir := filepath.Join(homeDir, ".config", "fish", "completions")
		fishCompletionFile := filepath.Join(fishCompletionDir, "mark.fish")
		_, err := os.Stat(fishCompletionFile)
		return err == nil
	}
	return false
}

// CheckFileForCompletionSource checks if a file sources mark completion
func CheckFileForCompletionSource(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if (strings.Contains(line, "~/.mark.bash") || strings.Contains(line, "~/.mark.zsh")) &&
			(strings.Contains(line, "source") || strings.Contains(line, ".")) ||
			(strings.Contains(line, "mark") && (strings.Contains(line, "complete") || strings.Contains(line, "completion"))) {
			return true
		}
	}
	return false
}

// SetupBashCompletion sets up bash command completion
func SetupBashCompletion() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		return
	}

	// Write the embedded completion script to ~/.mark.bash
	completionScriptPath := filepath.Join(homeDir, ".mark.bash")
	bashCompletionScript := `#!/bin/bash

# Helper function to get bookmarks with their paths for display
_mark_list_with_paths() {
    if [[ ! -d ~/.marks ]]; then
        return
    fi

    local mark target
    for mark in ~/.marks/*; do
        if [[ -L "$mark" ]]; then
            target=$(readlink "$mark" 2>/dev/null || echo "[broken]")
            printf "%-20s -> %s\n" "$(basename "$mark")" "$target"
        fi
    done | sort
}

_mark_complete() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"
    local cmd="${COMP_WORDS[0]}"

    # If we're on the first argument
    if [[ ${COMP_CWORD} -eq 1 ]]; then
        # If user starts typing a dash, offer flags (only for 'mark' command)
        if [[ "$cur" == -* && "$cmd" == "mark" ]]; then
            local flags="-l -d -j -v -h --config --autocomplete --alias --help --version"
            COMPREPLY=($(compgen -W "$flags" -- "${cur}"))
        else
            # For bookmark completion, show formatted list
            if [[ -d ~/.marks ]]; then
                # Get bookmark names for actual completion
                local marks=$(ls ~/.marks 2>/dev/null | tr '\n' ' ')
                COMPREPLY=($(compgen -W "$marks" -- "${cur}"))

                # If there are multiple matches or user hit tab twice, show formatted list
                if [[ ${#COMPREPLY[@]} -gt 1 ]]; then
                    echo >&2  # Newline before the list
                    _mark_list_with_paths >&2
                fi
            fi
        fi
    # If previous was -d or -j, offer bookmark names with paths
    elif [[ "$prev" == "-d" || "$prev" == "-j" ]]; then
        if [[ -d ~/.marks ]]; then
            local marks=$(ls ~/.marks 2>/dev/null | tr '\n' ' ')
            COMPREPLY=($(compgen -W "$marks" -- "${cur}"))

            # Show formatted list when multiple matches
            if [[ ${#COMPREPLY[@]} -gt 1 ]]; then
                echo >&2  # Newline before the list
                _mark_list_with_paths >&2
            fi
        fi
    fi
}

complete -F _mark_complete mark
complete -F _mark_complete marks
complete -F _mark_complete unmark
complete -F _mark_complete jump
`

	if err := os.WriteFile(completionScriptPath, []byte(bashCompletionScript), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing bash completion script: %v\n", err)
		return
	}

	// Add source line to .bashrc
	bashrc := filepath.Join(homeDir, ".bashrc")
	sourceLine := fmt.Sprintf("\n# mark command completion\nsource ~/.mark.bash\n")

	file, err := os.OpenFile(bashrc, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening .bashrc: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(sourceLine); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating .bashrc: %v\n", err)
		return
	}

	fmt.Printf("✓ Bash completion setup complete!\n")
	fmt.Printf("  Created completion script at %s\n", completionScriptPath)
	fmt.Printf("  Updated %s to source completion\n", bashrc)
	fmt.Printf("  Run 'source ~/.bashrc' or restart your shell to activate completions\n")
}

// SetupZshCompletion sets up zsh command completion
func SetupZshCompletion() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		return
	}

	// Write the embedded completion script to ~/.mark.zsh
	completionScriptPath := filepath.Join(homeDir, ".mark.zsh")
	zshCompletionScript := `#!/bin/zsh

_mark_complete() {
    local cur="${words[CURRENT]}"
    local prev="${words[CURRENT-1]}"
    local cmd="${words[1]}"

    # If we're on the first argument
    if [[ $CURRENT -eq 2 ]]; then
        # If user starts typing a dash, offer flags (only for 'mark' command)
        if [[ "$cur" == -* && "$cmd" == "mark" ]]; then
            local flags=("-l" "-d" "-j" "-v" "-h" "--config" "--autocomplete" "--alias" "--help" "--version")
            compadd -a flags
        else
            # For bookmark completion, offer with descriptions
            if [[ -d ~/.marks ]]; then
                local -a marks descriptions
                local mark target

                # Build arrays of marks and their descriptions
                for mark in ~/.marks/*(.N); do
                    if [[ -L "$mark" ]]; then
                        target=$(readlink "$mark" 2>/dev/null || echo "[broken]")
                        marks+=($(basename "$mark"))
                        descriptions+=("-> $target")
                    fi
                done

                # Use compadd with descriptions
                if [[ ${#marks[@]} -gt 0 ]]; then
                    compadd -d descriptions -a marks
                fi
            fi
        fi

    # If previous was -d or -j, offer bookmark names with descriptions
    elif [[ "$prev" == "-d" || "$prev" == "-j" ]]; then
        if [[ -d ~/.marks ]]; then
            local -a marks descriptions
            local mark target

            # Build arrays of marks and their descriptions
            for mark in ~/.marks/*(.N); do
                if [[ -L "$mark" ]]; then
                    target=$(readlink "$mark" 2>/dev/null || echo "[broken]")
                    marks+=($(basename "$mark"))
                    descriptions+=("-> $target")
                fi
            done

            # Use compadd with descriptions
            if [[ ${#marks[@]} -gt 0 ]]; then
                compadd -d descriptions -a marks
            fi
        fi
    fi
}

compdef _mark_complete mark
compdef _mark_complete marks
compdef _mark_complete unmark
compdef _mark_complete jump
`

	if err := os.WriteFile(completionScriptPath, []byte(zshCompletionScript), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing completion script: %v\n", err)
		return
	}

	// Add source line to .zshrc
	zshrcPath := filepath.Join(homeDir, ".zshrc")
	sourceLine := fmt.Sprintf("\n# mark command completion\nautoload -U +X compinit && compinit\nsource %s\n", completionScriptPath)

	file, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening .zshrc: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(sourceLine); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to .zshrc: %v\n", err)
		return
	}

	fmt.Printf("✓ Zsh completion setup complete!\n")
	fmt.Printf("  Created completion script at %s\n", completionScriptPath)
	fmt.Printf("  Added source line to %s\n", zshrcPath)
	fmt.Printf("  Restart your shell or run: source %s\n", zshrcPath)
}

// SetupFishCompletion sets up fish command completion
func SetupFishCompletion() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		return
	}

	// Create fish completion directory if it doesn't exist
	fishCompletionDir := filepath.Join(homeDir, ".config", "fish", "completions")
	if err := os.MkdirAll(fishCompletionDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating fish completion directory: %v\n", err)
		return
	}

	// Create a fish completion script
	fishCompletionScript := `# mark command completion for fish

# Helper function to list bookmarks with their paths
function __fish_mark_list_bookmarks
    if test -d ~/.marks
        for mark in ~/.marks/*
            if test -L "$mark"
                set -l target (readlink "$mark" 2>/dev/null; or echo "[broken]")
                set -l name (basename "$mark")
                echo -e "$name\t-> $target"
            end
        end
    end
end

complete -c mark -f
complete -c mark -s l -d "List bookmarks"
complete -c mark -s d -d "Delete bookmark" -r
complete -c mark -s j -d "Jump to bookmark" -r
complete -c mark -l config -d "Run setup/reconfigure"
complete -c mark -l autocomplete -d "Setup/update command line autocompletion"
complete -c mark -l alias -d "Setup shell aliases"
complete -c mark -s v -l version -d "Show version"
complete -c mark -s h -l help -d "Show help"

# Complete with existing bookmark names with paths for main argument
complete -c mark -n '__fish_is_first_token' -a '(__fish_mark_list_bookmarks)'

# Complete with bookmark names and paths for -d and -j flags
complete -c mark -n '__fish_seen_subcommand_from -d' -a '(__fish_mark_list_bookmarks)'
complete -c mark -n '__fish_seen_subcommand_from -j' -a '(__fish_mark_list_bookmarks)'

# Alias completions with descriptions
complete -c marks -f -a '(__fish_mark_list_bookmarks)'
complete -c unmark -f -a '(__fish_mark_list_bookmarks)'
complete -c jump -f -a '(__fish_mark_list_bookmarks)'
`

	markCompletionFile := filepath.Join(fishCompletionDir, "mark.fish")
	if err := os.WriteFile(markCompletionFile, []byte(fishCompletionScript), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing fish completion script: %v\n", err)
		return
	}

	fmt.Printf("✓ Fish completion setup complete!\n")
	fmt.Printf("  Created completion file at %s\n", markCompletionFile)
	fmt.Printf("  Restart your shell to activate completions\n")
}

// RunAutocompleteSetup handles the main autocomplete setup flow
func RunAutocompleteSetup() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("mark - Command Line Autocompletion Setup")
	fmt.Println()
	fmt.Println("This will set up tab completion for the mark command, allowing you to:")
	fmt.Println("• Tab-complete bookmark names")
	fmt.Println("• Tab-complete command flags")
	fmt.Println("• Get context-aware completions")
	fmt.Println()
	fmt.Print("Would you like to set up autocompletion? (y/N): ")

	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "y" && response != "yes" {
		fmt.Println("Autocompletion setup cancelled.")
		return
	}

	shell := detectShell()
	if shell == "" {
		fmt.Println("Could not detect shell type. Skipping completion setup.")
		fmt.Println("Supported shells: bash, zsh, fish")
		return
	}

	fmt.Printf("Detected shell: %s\n", shell)
	fmt.Println()

	// Clean up any existing completion setup
	fmt.Println("Cleaning up any existing completion setup...")
	CleanupExistingCompletion(shell)

	// Set up completion for the detected shell
	fmt.Printf("Setting up %s completion...\n", shell)
	switch shell {
	case "bash":
		SetupBashCompletion()
	case "zsh":
		SetupZshCompletion()
	case "fish":
		SetupFishCompletion()
	default:
		fmt.Printf("Shell '%s' not supported for completion. Supported shells: bash, zsh, fish\n", shell)
		return
	}

	fmt.Println()
	fmt.Println("✓ Autocompletion setup complete!")
	fmt.Println("  To activate, run one of:")

	homeDir, _ := os.UserHomeDir()
	switch shell {
	case "bash":
		fmt.Printf("    source ~/.bashrc\n")
		fmt.Printf("    source %s\n", filepath.Join(homeDir, ".mark.bash"))
	case "zsh":
		fmt.Printf("    source ~/.zshrc\n")
		fmt.Printf("    source %s\n", filepath.Join(homeDir, ".mark.zsh"))
	case "fish":
		fmt.Println("    (restart your shell)")
	}
	fmt.Println("  Or simply restart your shell")
}

// CleanupExistingCompletion removes existing completion setup for the specified shell
func CleanupExistingCompletion(shell string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	switch shell {
	case "bash":
		// Remove existing .mark.bash file
		bashCompletionFile := filepath.Join(homeDir, ".mark.bash")
		os.Remove(bashCompletionFile)

		// Clean up shell config files
		cleanupShellConfig(filepath.Join(homeDir, ".bashrc"))
		cleanupShellConfig(filepath.Join(homeDir, ".bash_profile"))
		cleanupShellConfig(filepath.Join(homeDir, ".profile"))

	case "zsh":
		// Remove existing .mark.zsh file
		zshCompletionFile := filepath.Join(homeDir, ".mark.zsh")
		os.Remove(zshCompletionFile)

		// Clean up .zshrc
		cleanupShellConfig(filepath.Join(homeDir, ".zshrc"))

	case "fish":
		// Remove existing fish completion file
		fishCompletionDir := filepath.Join(homeDir, ".config", "fish", "completions")
		markCompletionFile := filepath.Join(fishCompletionDir, "mark.fish")
		os.Remove(markCompletionFile)
	}
}

// cleanupShellConfig removes mark completion lines from shell config files
func cleanupShellConfig(configFile string) {
	// Read the file
	file, err := os.Open(configFile)
	if err != nil {
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	skipNext := false

	for scanner.Scan() {
		line := scanner.Text()

		// Skip lines that contain mark completion references
		if strings.Contains(line, "# mark command completion") {
			skipNext = true
			continue
		}

		if skipNext && (strings.Contains(line, ".mark.bash") ||
			strings.Contains(line, ".mark.zsh") ||
			strings.Contains(line, "completions/bash/mark") ||
			(strings.Contains(line, "mark") && strings.Contains(line, "source"))) {
			skipNext = false
			continue
		}

		if skipNext && strings.TrimSpace(line) == "" {
			continue
		}

		skipNext = false
		lines = append(lines, line)
	}

	// Write the cleaned file back
	outFile, err := os.Create(configFile)
	if err != nil {
		return
	}
	defer outFile.Close()

	for _, line := range lines {
		fmt.Fprintln(outFile, line)
	}
}

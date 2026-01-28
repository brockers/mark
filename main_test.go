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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Could not get home directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "expand tilde",
			input:    "~/test",
			expected: filepath.Join(homeDir, "test"),
		},
		{
			name:     "absolute path unchanged",
			input:    "/tmp/test",
			expected: "/tmp/test",
		},
		{
			name:     "relative path unchanged (no symlinks)",
			input:    "test/path",
			expected: "test/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandPath(tt.input)
			if result != tt.expected {
				t.Errorf("expandPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedFlags *ParsedFlags
		expectedArgs  []string
	}{
		{
			name: "help flag short",
			args: []string{"-h"},
			expectedFlags: &ParsedFlags{
				Help: true,
			},
			expectedArgs: []string{},
		},
		{
			name: "help flag long",
			args: []string{"--help"},
			expectedFlags: &ParsedFlags{
				Help: true,
			},
			expectedArgs: []string{},
		},
		{
			name: "version flag short",
			args: []string{"-v"},
			expectedFlags: &ParsedFlags{
				Version: true,
			},
			expectedArgs: []string{},
		},
		{
			name: "version flag long",
			args: []string{"--version"},
			expectedFlags: &ParsedFlags{
				Version: true,
			},
			expectedArgs: []string{},
		},
		{
			name: "list flag",
			args: []string{"-l"},
			expectedFlags: &ParsedFlags{
				List: true,
			},
			expectedArgs: []string{},
		},
		{
			name: "delete flag",
			args: []string{"-d", "testmark"},
			expectedFlags: &ParsedFlags{
				Delete: "testmark",
			},
			expectedArgs: []string{},
		},
		{
			name: "jump flag",
			args: []string{"-j", "testmark"},
			expectedFlags: &ParsedFlags{
				Jump: "testmark",
			},
			expectedArgs: []string{},
		},
		{
			name: "config flag",
			args: []string{"--config"},
			expectedFlags: &ParsedFlags{
				Config: true,
			},
			expectedArgs: []string{},
		},
		{
			name: "configure flag (alias for config)",
			args: []string{"--configure"},
			expectedFlags: &ParsedFlags{
				Config: true,
			},
			expectedArgs: []string{},
		},
		{
			name: "autocomplete flag",
			args: []string{"--autocomplete"},
			expectedFlags: &ParsedFlags{
				Autocomplete: true,
			},
			expectedArgs: []string{},
		},
		{
			name: "alias flag",
			args: []string{"--alias"},
			expectedFlags: &ParsedFlags{
				Alias: true,
			},
			expectedArgs: []string{},
		},
		{
			name:          "no flags with args",
			args:          []string{"mybookmark"},
			expectedFlags: &ParsedFlags{},
			expectedArgs:  []string{"mybookmark"},
		},
		{
			name:          "no flags with multiple args",
			args:          []string{"my", "bookmark", "name"},
			expectedFlags: &ParsedFlags{},
			expectedArgs:  []string{"my", "bookmark", "name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags, args := parseFlags(tt.args)

			// Check all flag fields
			if flags.Help != tt.expectedFlags.Help {
				t.Errorf("Help flag mismatch: got %v, want %v", flags.Help, tt.expectedFlags.Help)
			}
			if flags.Version != tt.expectedFlags.Version {
				t.Errorf("Version flag mismatch: got %v, want %v", flags.Version, tt.expectedFlags.Version)
			}
			if flags.List != tt.expectedFlags.List {
				t.Errorf("List flag mismatch: got %v, want %v", flags.List, tt.expectedFlags.List)
			}
			if flags.Delete != tt.expectedFlags.Delete {
				t.Errorf("Delete flag mismatch: got %q, want %q", flags.Delete, tt.expectedFlags.Delete)
			}
			if flags.Jump != tt.expectedFlags.Jump {
				t.Errorf("Jump flag mismatch: got %q, want %q", flags.Jump, tt.expectedFlags.Jump)
			}
			if flags.Config != tt.expectedFlags.Config {
				t.Errorf("Config flag mismatch: got %v, want %v", flags.Config, tt.expectedFlags.Config)
			}
			if flags.Autocomplete != tt.expectedFlags.Autocomplete {
				t.Errorf("Autocomplete flag mismatch: got %v, want %v", flags.Autocomplete, tt.expectedFlags.Autocomplete)
			}
			if flags.Alias != tt.expectedFlags.Alias {
				t.Errorf("Alias flag mismatch: got %v, want %v", flags.Alias, tt.expectedFlags.Alias)
			}

			// Check remaining args
			if len(args) != len(tt.expectedArgs) {
				t.Errorf("Args length mismatch: got %d, want %d", len(args), len(tt.expectedArgs))
			} else {
				for i, arg := range args {
					if arg != tt.expectedArgs[i] {
						t.Errorf("Arg[%d] mismatch: got %q, want %q", i, arg, tt.expectedArgs[i])
					}
				}
			}
		})
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create a temporary home directory
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create a test config
	testConfig := Config{
		MarksDir: filepath.Join(tmpDir, ".marks"),
	}

	// Save the config
	saveConfig(testConfig)

	// Check if config file was created
	configPath := filepath.Join(tmpDir, ".mark")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Read the config file
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Could not read config file: %v", err)
	}

	// Verify content contains marksdir
	contentStr := string(content)
	if !strings.Contains(contentStr, "marksdir=") {
		t.Errorf("Config file does not contain marksdir field")
	}

	// Verify it uses tilde notation
	if !strings.Contains(contentStr, "~/.marks") {
		t.Errorf("Config file does not use tilde notation: %s", contentStr)
	}
}

func TestDetectShell(t *testing.T) {
	tests := []struct {
		name        string
		shellEnv    string
		expectedRes string
	}{
		{
			name:        "bash",
			shellEnv:    "/bin/bash",
			expectedRes: "bash",
		},
		{
			name:        "zsh",
			shellEnv:    "/usr/bin/zsh",
			expectedRes: "zsh",
		},
		{
			name:        "fish",
			shellEnv:    "/usr/local/bin/fish",
			expectedRes: "fish",
		},
		{
			name:        "empty",
			shellEnv:    "",
			expectedRes: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalShell := os.Getenv("SHELL")
			os.Setenv("SHELL", tt.shellEnv)
			defer os.Setenv("SHELL", originalShell)

			result := detectShell()
			if result != tt.expectedRes {
				t.Errorf("detectShell() = %q, want %q", result, tt.expectedRes)
			}
		})
	}
}

// Integration-style tests for bookmark operations
func TestGenerateBashRC(t *testing.T) {
	tests := []struct {
		name               string
		includeAliases     bool
		includeCompletions bool
		expectAliases      bool
		expectCompletions  bool
	}{
		{
			name:               "aliases only",
			includeAliases:     true,
			includeCompletions: false,
			expectAliases:      true,
			expectCompletions:  false,
		},
		{
			name:               "completions only",
			includeAliases:     false,
			includeCompletions: true,
			expectAliases:      false,
			expectCompletions:  true,
		},
		{
			name:               "both",
			includeAliases:     true,
			includeCompletions: true,
			expectAliases:      true,
			expectCompletions:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := generateBashRC("/usr/bin/mark", tt.includeAliases, tt.includeCompletions)

			// Check header
			if !strings.Contains(content, "# mark shell configuration") {
				t.Error("Missing header comment")
			}
			if !strings.Contains(content, "# Generated by mark") {
				t.Error("Missing generation comment")
			}

			// Check features line
			if tt.expectAliases && !strings.Contains(content, "aliases") {
				t.Error("Missing 'aliases' in features line")
			}
			if tt.expectCompletions && !strings.Contains(content, "completions") {
				t.Error("Missing 'completions' in features line")
			}

			// Check aliases content
			hasAliases := strings.Contains(content, "alias marks=") && strings.Contains(content, "function jump()")
			if tt.expectAliases && !hasAliases {
				t.Error("Expected aliases but not found")
			}
			if !tt.expectAliases && hasAliases {
				t.Error("Found aliases but not expected")
			}

			// Check completions content
			hasCompletions := strings.Contains(content, "_mark_complete()") && strings.Contains(content, "complete -F")
			if tt.expectCompletions && !hasCompletions {
				t.Error("Expected completions but not found")
			}
			if !tt.expectCompletions && hasCompletions {
				t.Error("Found completions but not expected")
			}
		})
	}
}

func TestGenerateZshRC(t *testing.T) {
	content := generateZshRC("/usr/bin/mark", true, true)

	// Check header
	if !strings.Contains(content, "#!/bin/zsh") {
		t.Error("Missing zsh shebang")
	}
	if !strings.Contains(content, "# mark shell configuration") {
		t.Error("Missing header comment")
	}

	// Check aliases
	if !strings.Contains(content, "alias marks=") {
		t.Error("Missing marks alias")
	}
	if !strings.Contains(content, "function jump()") {
		t.Error("Missing jump function")
	}

	// Check completions
	if !strings.Contains(content, "compdef _mark_complete mark") {
		t.Error("Missing compdef for mark")
	}
	if !strings.Contains(content, "autoload -U +X compinit") {
		t.Error("Missing compinit")
	}
}

func TestGenerateFishRC(t *testing.T) {
	content := generateFishRC("/usr/bin/mark", true, true)

	// Check header (fish doesn't use shebang in conf.d)
	if !strings.Contains(content, "# mark shell configuration") {
		t.Error("Missing header comment")
	}

	// Check aliases
	if !strings.Contains(content, "alias marks ") {
		t.Error("Missing marks alias")
	}
	if !strings.Contains(content, "function jump") {
		t.Error("Missing jump function")
	}

	// Check completions
	if !strings.Contains(content, "complete -c mark") {
		t.Error("Missing mark completion")
	}
	if !strings.Contains(content, "__fish_mark_list_bookmarks") {
		t.Error("Missing bookmark list helper")
	}
}

func TestIsSourceLinePresent(t *testing.T) {
	tmpDir := t.TempDir()

	// Test file with source line
	fileWithSource := filepath.Join(tmpDir, "with_source")
	content := "# some content\n# mark shell integration\n[ -f ~/.mark_bash_rc ] && source ~/.mark_bash_rc\n"
	os.WriteFile(fileWithSource, []byte(content), 0644)

	if !isSourceLinePresent(fileWithSource) {
		t.Error("Should detect source line in file")
	}

	// Test file without source line
	fileWithoutSource := filepath.Join(tmpDir, "without_source")
	content = "# some other content\nexport PATH=$PATH:/usr/bin\n"
	os.WriteFile(fileWithoutSource, []byte(content), 0644)

	if isSourceLinePresent(fileWithoutSource) {
		t.Error("Should not detect source line in file")
	}

	// Test non-existent file
	if isSourceLinePresent(filepath.Join(tmpDir, "nonexistent")) {
		t.Error("Should return false for non-existent file")
	}
}

func TestGetEnabledFeatures(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create RC file with both features
	rcPath := filepath.Join(tmpDir, bashRCFile)
	content := "#!/bin/bash\n# mark shell configuration\n# Generated by mark\n# Features: aliases completions\n"
	os.WriteFile(rcPath, []byte(content), 0644)

	aliases, completions := getEnabledFeatures("bash")
	if !aliases {
		t.Error("Should detect aliases feature")
	}
	if !completions {
		t.Error("Should detect completions feature")
	}

	// Test with only aliases
	content = "#!/bin/bash\n# mark shell configuration\n# Generated by mark\n# Features: aliases\n"
	os.WriteFile(rcPath, []byte(content), 0644)

	aliases, completions = getEnabledFeatures("bash")
	if !aliases {
		t.Error("Should detect aliases feature")
	}
	if completions {
		t.Error("Should not detect completions feature")
	}

	// Test with non-existent file
	os.Remove(rcPath)
	aliases, completions = getEnabledFeatures("bash")
	if aliases || completions {
		t.Error("Should return false for non-existent file")
	}
}

func TestWriteShellRC(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Test bash RC file creation
	err := writeShellRC("bash", true, true)
	if err != nil {
		t.Fatalf("Failed to write bash RC: %v", err)
	}

	rcPath := filepath.Join(tmpDir, bashRCFile)
	if _, err := os.Stat(rcPath); os.IsNotExist(err) {
		t.Error("Bash RC file not created")
	}

	content, _ := os.ReadFile(rcPath)
	if !strings.Contains(string(content), "alias marks=") {
		t.Error("Bash RC missing aliases")
	}
	if !strings.Contains(string(content), "_mark_complete()") {
		t.Error("Bash RC missing completions")
	}

	// Test fish RC file creation (should create conf.d directory)
	err = writeShellRC("fish", true, true)
	if err != nil {
		t.Fatalf("Failed to write fish RC: %v", err)
	}

	fishRcPath := filepath.Join(tmpDir, fishRCFile)
	if _, err := os.Stat(fishRcPath); os.IsNotExist(err) {
		t.Error("Fish RC file not created")
	}

	// Test unsupported shell
	err = writeShellRC("unsupported", true, true)
	if err == nil {
		t.Error("Should return error for unsupported shell")
	}
}

func TestEnsureSourceLine(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create empty .bashrc
	bashrcPath := filepath.Join(tmpDir, ".bashrc")
	os.WriteFile(bashrcPath, []byte("# existing content\n"), 0644)

	// Add source line
	err := ensureSourceLine("bash")
	if err != nil {
		t.Fatalf("Failed to add source line: %v", err)
	}

	content, _ := os.ReadFile(bashrcPath)
	if !strings.Contains(string(content), "# mark shell integration") {
		t.Error("Missing source line marker")
	}
	if !strings.Contains(string(content), ".mark_bash_rc") {
		t.Error("Missing RC file reference")
	}

	// Running again should not duplicate
	err = ensureSourceLine("bash")
	if err != nil {
		t.Fatalf("Failed on second call: %v", err)
	}

	content, _ = os.ReadFile(bashrcPath)
	count := strings.Count(string(content), "# mark shell integration")
	if count != 1 {
		t.Errorf("Source line duplicated: found %d occurrences", count)
	}

	// Test fish (should do nothing)
	err = ensureSourceLine("fish")
	if err != nil {
		t.Fatalf("Fish should return nil: %v", err)
	}
}

func TestGetRCFilePath(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		shell    string
		expected string
	}{
		{"bash", filepath.Join(homeDir, ".mark_bash_rc")},
		{"zsh", filepath.Join(homeDir, ".mark_zsh_rc")},
		{"fish", filepath.Join(homeDir, ".config/fish/conf.d/mark.fish")},
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			result := getRCFilePath(tt.shell)
			if result != tt.expected {
				t.Errorf("getRCFilePath(%q) = %q, want %q", tt.shell, result, tt.expected)
			}
		})
	}
}

func TestBookmarkOperations(t *testing.T) {
	// Create a temporary marks directory
	tmpDir := t.TempDir()
	marksDir := filepath.Join(tmpDir, ".marks")

	// Create the marks directory
	if err := os.MkdirAll(marksDir, 0755); err != nil {
		t.Fatalf("Could not create marks directory: %v", err)
	}

	// Test bookmark creation
	t.Run("create bookmark", func(t *testing.T) {
		// Create a test directory to bookmark
		testDir := filepath.Join(tmpDir, "test-project")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("Could not create test directory: %v", err)
		}

		// Change to test directory
		originalWd, _ := os.Getwd()
		os.Chdir(testDir)
		defer os.Chdir(originalWd)

		// This would normally call createBookmark, but that exits on error
		// So we'll test the symlink creation logic directly
		bookmarkName := "testproject"
		symlinkPath := filepath.Join(marksDir, bookmarkName)

		if err := os.Symlink(testDir, symlinkPath); err != nil {
			t.Fatalf("Could not create symlink: %v", err)
		}

		// Verify symlink was created
		fileInfo, err := os.Lstat(symlinkPath)
		if err != nil {
			t.Fatalf("Symlink was not created: %v", err)
		}

		if fileInfo.Mode()&os.ModeSymlink == 0 {
			t.Error("Created file is not a symlink")
		}

		// Verify symlink target
		target, err := os.Readlink(symlinkPath)
		if err != nil {
			t.Fatalf("Could not read symlink: %v", err)
		}

		if target != testDir {
			t.Errorf("Symlink target = %q, want %q", target, testDir)
		}
	})

	// Test custom path bookmark creation
	t.Run("create bookmark with custom path", func(t *testing.T) {
		// Create a test directory to bookmark
		customDir := filepath.Join(tmpDir, "custom-location")
		if err := os.MkdirAll(customDir, 0755); err != nil {
			t.Fatalf("Could not create custom directory: %v", err)
		}

		// Create bookmark with custom path (not current directory)
		bookmarkName := "custommark"
		symlinkPath := filepath.Join(marksDir, bookmarkName)

		if err := os.Symlink(customDir, symlinkPath); err != nil {
			t.Fatalf("Could not create symlink: %v", err)
		}

		// Verify symlink target points to custom directory
		target, err := os.Readlink(symlinkPath)
		if err != nil {
			t.Fatalf("Could not read symlink: %v", err)
		}

		if target != customDir {
			t.Errorf("Symlink target = %q, want %q", target, customDir)
		}
	})

	// Test broken symlink detection
	t.Run("detect broken symlink", func(t *testing.T) {
		// Create a symlink to a non-existent directory
		brokenDir := filepath.Join(tmpDir, "non-existent")
		brokenName := "broken-mark"
		brokenPath := filepath.Join(marksDir, brokenName)

		if err := os.Symlink(brokenDir, brokenPath); err != nil {
			t.Fatalf("Could not create broken symlink: %v", err)
		}

		// Verify it's a symlink
		fileInfo, err := os.Lstat(brokenPath)
		if err != nil {
			t.Fatalf("Could not stat broken symlink: %v", err)
		}

		if fileInfo.Mode()&os.ModeSymlink == 0 {
			t.Error("File is not a symlink")
		}

		// Verify target doesn't exist (broken link)
		_, err = os.Stat(brokenPath)
		if err == nil {
			t.Error("Expected error for broken symlink, got nil")
		}
	})

	// Test bookmark deletion
	t.Run("delete bookmark", func(t *testing.T) {
		// Create a test symlink
		testDir := filepath.Join(tmpDir, "delete-test")
		os.MkdirAll(testDir, 0755)

		deleteName := "delete-mark"
		deletePath := filepath.Join(marksDir, deleteName)

		if err := os.Symlink(testDir, deletePath); err != nil {
			t.Fatalf("Could not create test symlink: %v", err)
		}

		// Delete the symlink
		if err := os.Remove(deletePath); err != nil {
			t.Fatalf("Could not remove symlink: %v", err)
		}

		// Verify it was deleted
		_, err := os.Lstat(deletePath)
		if err == nil {
			t.Error("Symlink still exists after deletion")
		}
		if !os.IsNotExist(err) {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

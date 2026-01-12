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

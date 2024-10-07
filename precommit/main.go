package precommit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetGitRoot returns the absolute path to the Git root directory.
func GetGitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// AddPreCommitHook creates a pre-commit hook in the Git repository.
func AddPreCommitHook() error {
	// Get the Git root directory
	gitRoot, err := GetGitRoot()
	if err != nil {
		return fmt.Errorf("failed to get Git root: %w", err)
	}

	// Define the pre-commit hook path
	hookPath := filepath.Join(gitRoot, ".git", "hooks", "commit-msg")

	// Define the content of the commit-msg hook script
	hookContent := fmt.Sprintf(`#!/bin/bash
# commit-msg hook to run anto validate script
INPUT_FILE=$1
# Navigate to the .anto directory
cd %s/.anto

# Explicitly run anto script with 'validate' as a parameter
bash -c './anto validate $INPUT_FILE'
`, gitRoot)

	// Create or overwrite the commit-msg hook file
	if err := os.WriteFile(hookPath, []byte(hookContent), 0755); err != nil {
		return fmt.Errorf("failed to write pre-commit hook: %w", err)
	}

	fmt.Println("Commit-msg hook created successfully at", hookPath)
	return nil
}

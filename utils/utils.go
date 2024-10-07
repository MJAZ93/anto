package utils

import (
	"anto/model"
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func GetProjectFileName() (string, error) {
	return GetGitRoot()
}

// GetGitRoot returns the absolute path to the Git root directory.
func GetGitRoot() (string, error) {
	// Command to find the git root directory
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	// Execute the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Clean up the output by removing any trailing newlines
	gitRoot := strings.TrimSpace(string(output))
	return gitRoot, nil
}

// GetGitRootFolderName returns the name of the Git root folder.
func GetGitRootFolderName() (string, error) {
	// Command to find the git root directory
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	// Execute the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Clean up the output by removing any trailing newlines
	gitRoot := strings.TrimSpace(string(output))

	// Extract the folder name from the Git root path
	folderName := filepath.Base(gitRoot)

	return folderName, nil
}

func ParseMsk(filename string) (*model.MskMessage, error) {
	cm := &model.MskMessage{
		// Initialize MaxLines to the maximum integer value
		MaxLines: math.MaxInt32,
		MinLines: 0,
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inBlockComment := false
	var commentLines []string

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines unless inside a block comment
		if trimmedLine == "" && !inBlockComment {
			continue
		}

		// Handle block comments
		if strings.Contains(trimmedLine, "/*") {
			inBlockComment = true
			// Extract content after /*
			startIndex := strings.Index(trimmedLine, "/*") + 2
			content := trimmedLine[startIndex:]

			// Check if the block comment ends on the same line
			if strings.Contains(content, "*/") {
				endIndex := strings.Index(content, "*/")
				commentLines = append(commentLines, strings.TrimSpace(content[:endIndex]))
				inBlockComment = false
			} else {
				commentLines = append(commentLines, strings.TrimSpace(content))
			}
			continue
		}

		if inBlockComment {
			// Check if the line contains */
			if strings.Contains(trimmedLine, "*/") {
				endIndex := strings.Index(trimmedLine, "*/")
				commentLines = append(commentLines, strings.TrimSpace(trimmedLine[:endIndex]))
				inBlockComment = false
			} else {
				commentLines = append(commentLines, trimmedLine)
			}
			continue
		}

		// Skip single-line comments
		if strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		// Process the line
		switch {
		case strings.HasPrefix(trimmedLine, "l"):
			parts := strings.Fields(trimmedLine)
			if len(parts) >= 2 && parts[0] == "l" {
				// Handle range with '..'
				if strings.Contains(parts[1], "..") {
					rangeParts := strings.Split(parts[1], "..")
					if len(rangeParts) != 2 {
						return nil, fmt.Errorf("invalid range in line: %s", line)
					}
					minLines, err1 := strconv.Atoi(rangeParts[0])
					maxLines, err2 := strconv.Atoi(rangeParts[1])
					if err1 != nil || err2 != nil {
						return nil, fmt.Errorf("invalid numbers in range: %s", line)
					}
					cm.MinLines = minLines
					cm.MaxLines = maxLines
				} else if len(parts) >= 3 {
					// Handle operators <, >, =
					n, err := strconv.Atoi(parts[1])
					if err != nil {
						return nil, fmt.Errorf("invalid number in line: %s", line)
					}
					operator := parts[2]
					switch operator {
					case "<":
						cm.MinLines = 0
						cm.MaxLines = n - 1
					case ">":
						cm.MinLines = n + 1
						cm.MaxLines = math.MaxInt32
					case "=":
						cm.MinLines = n
						cm.MaxLines = n
					default:
						return nil, fmt.Errorf("invalid operator in line: %s", line)
					}
				} else {
					return nil, fmt.Errorf("invalid line format: %s", line)
				}
			}
		case strings.HasPrefix(trimmedLine, "+"):
			rule := strings.TrimSpace(trimmedLine[1:])
			cm.AllowedRules = append(cm.AllowedRules, rule)
		case strings.HasPrefix(trimmedLine, "-"):
			rule := strings.TrimSpace(trimmedLine[1:])
			cm.DisallowedRules = append(cm.DisallowedRules, rule)
		default:
			// Ignore any other lines
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Combine collected comment lines into a single string
	cm.Comment = strings.Join(commentLines, "\n")

	return cm, nil
}

func ComparePaths(path1, path2 string) (samePart, diffPart1, diffPart2 string) {
	parts1 := strings.Split(path1, "/")
	parts2 := strings.Split(path2, "/")

	var sameParts []string
	var diffParts1 []string
	var diffParts2 []string

	// Find the common parts
	minLen := len(parts1)
	if len(parts2) < minLen {
		minLen = len(parts2)
	}

	for i := 0; i < minLen; i++ {
		if parts1[i] == parts2[i] {
			sameParts = append(sameParts, parts1[i])
		} else {
			diffParts1 = append(diffParts1, parts1[i])
			diffParts2 = append(diffParts2, parts2[i])
		}
	}

	// Append remaining parts if the paths are of different lengths
	if len(parts1) > minLen {
		diffParts1 = append(diffParts1, parts1[minLen:]...)
	}
	if len(parts2) > minLen {
		diffParts2 = append(diffParts2, parts2[minLen:]...)
	}

	// Convert the parts back to strings
	samePart = strings.Join(sameParts, "/")
	diffPart1 = strings.Join(diffParts1, "/")
	diffPart2 = strings.Join(diffParts2, "/")

	return
}

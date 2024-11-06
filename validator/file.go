package validator

import (
	"anto/generator"
	"anto/model"
	"anto/utils"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func ValidateStructure(folder model.FolderRule, fileRoot string) error {
	// Get list of folders and files in the directory
	folders, files, err := getFsFsInDirectory(fileRoot)
	if err != nil {
		return fmt.Errorf("error listing folders and files in %s: %v", fileRoot, err)
	}

	usedFolders := make(map[string]bool)
	usedFiles := make(map[string]bool)

	// Flags to control ignoring files and folders
	ignoreAllFolders := false
	ignoreAllFiles := false

	// Handle special folder patterns
	for _, folderRule := range folder.Folders {
		if folderRule.Pattern == "**" {
			// Ignore all files and folders in this directory and subdirectories
			return nil // Stop validation in this folder
		} else if folderRule.Pattern == "*" {
			// Ignore all folders in the current directory
			ignoreAllFolders = true
			// Mark all folders as used
			for _, f := range folders {
				usedFolders[f] = true
			}
			break // No need to process other folder rules
		}
	}

	// Handle special file patterns
	for _, fileRule := range folder.Files {
		if fileRule.Pattern == "**" {
			// Ignore all files in this directory and subdirectories
			ignoreAllFiles = true
			// Mark all files as used
			for _, f := range files {
				usedFiles[f] = true
			}
			break // No need to process other file rules
		} else if fileRule.Pattern == "*" {
			// Ignore all files in the current directory
			ignoreAllFiles = true
			// Mark all files as used
			for _, f := range files {
				usedFiles[f] = true
			}
			break // No need to process other file rules
		}
	}

	// If not ignoring all folders, proceed to validate folders
	if !ignoreAllFolders {
		// Validate folders
		for _, folderRule := range folder.Folders {
			// Skip special patterns as we've already handled them
			if folderRule.Pattern == "*" || folderRule.Pattern == "**" {
				continue
			}

			matchingFolders := []string{}
			for _, f := range folders {
				matches, err := matchesPattern(f, folderRule.Pattern)
				if err != nil {
					return fmt.Errorf("invalid pattern %s: %v", folderRule.Pattern, err)
				}
				if matches {
					matchingFolders = append(matchingFolders, f)
					usedFolders[f] = true
				}
			}
			count := len(matchingFolders)
			if count < folderRule.MinCount || count > folderRule.MaxCount {
				return fmt.Errorf("in %s: expected between %d and %d folders matching pattern '%s', found %d",
					fileRoot, folderRule.MinCount, folderRule.MaxCount, folderRule.Pattern, count)
			}
			for _, matchedFolder := range matchingFolders {
				err := ValidateStructure(folderRule, filepath.Join(fileRoot, matchedFolder))
				if err != nil {
					return err
				}
			}
		}
	}

	// If not ignoring all files, proceed to validate files
	if !ignoreAllFiles {
		// Validate files
		for _, fileRule := range folder.Files {
			// Skip special patterns as we've already handled them
			if fileRule.Pattern == "*" || fileRule.Pattern == "**" {
				continue
			}

			matchingFiles := []string{}
			for _, f := range files {
				matches, err := matchesPattern(f, fileRule.Pattern)
				if err != nil {
					return fmt.Errorf("invalid pattern %s: %v", fileRule.Pattern, err)
				}
				if matches {
					matchingFiles = append(matchingFiles, f)
					usedFiles[f] = true
				}

				err = checksFileRegex(fileRoot, f, fileRule.Pattern)
				if err != nil {
					return err
				}
			}
			count := len(matchingFiles)
			if count < fileRule.MinCount || count > fileRule.MaxCount {
				return fmt.Errorf("in %s: expected between %d and %d files matching pattern '%s', found %d",
					fileRoot, fileRule.MinCount, fileRule.MaxCount, fileRule.Pattern, count)
			}
		}
	}

	// Check for unexpected folders (if not ignoring all folders)
	if !ignoreAllFolders {
		for _, f := range folders {
			if !usedFolders[f] && fileRoot != "./" && !strings.HasPrefix(f, ".") {
				return fmt.Errorf("unexpected folder '%s' in %s", f, fileRoot)
			}
		}
	}

	// Check for unexpected files (if not ignoring all files)
	if !ignoreAllFiles {
		for _, f := range files {
			if !usedFiles[f] && fileRoot != "./" && !strings.HasPrefix(f, ".") {
				return fmt.Errorf("unexpected file '%s' in %s", f, fileRoot)
			}
		}
	}

	return nil
}

// Rest of your helper functions remain the same...

func checksFileRegex(fileRoot string, filename string, regex string) error {
	gitRoot, err := utils.GetGitRoot()
	if err != nil {
		return err
	}

	gitFolder, err := utils.GetGitRootFolderName()
	if err != nil {
		return err
	}

	samePart, _, diffPart2 := utils.ComparePaths(gitRoot, fileRoot)

	if !strings.HasPrefix("/", diffPart2) {
		diffPart2 = "/" + diffPart2
	}

	mksFile := samePart + "/.anto/" + gitFolder + diffPart2 + "/" + generator.SanitizeName(regex) + ".msk"

	actualFile := filepath.Join(fileRoot, filename)

	if fileExists(mksFile) {

		fileByte, err := os.ReadFile(actualFile)
		if err != nil {
			return errors.New("error reading file " + actualFile + ", err: " + err.Error())
		}
		fileContent := string(fileByte)

		cm, err := utils.ParseMsk(mksFile)
		if err != nil {
			return errors.New("error parsing " + mksFile + ", err: " + err.Error())
		}

		lines := strings.Split(fileContent, "\n")
		lineCount := len(lines)

		if lineCount > cm.MaxLines {
			return errors.New(actualFile + " has " + strconv.Itoa(lineCount) + " lines, when max allowed is " + strconv.Itoa(cm.MaxLines))
		}

		if lineCount < cm.MinLines {
			return errors.New(actualFile + " has " + strconv.Itoa(lineCount) + " lines, when min allowed is " + strconv.Itoa(cm.MinLines))
		}

		matchesAnyRule := false
		for _, rule := range cm.AllowedRules {
			re := regexp.MustCompile(rule)
			if re.MatchString(fileContent) {
				matchesAnyRule = true
				break
			}
		}

		if !matchesAnyRule && len(cm.AllowedRules) > 0 {
			return errors.New(actualFile + " doesn't match any allowed rule, file instructions: \n\n" + cm.Comment)
		}

		for _, rule := range cm.DisallowedRules {
			re := regexp.MustCompile(rule)
			if re.MatchString(fileContent) {
				return errors.New(actualFile + " matches disallowed rule: " + rule + ". File instructions: \n\n" + cm.Comment)
			}
		}

		return nil
	}

	return nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || !os.IsNotExist(err)
}

func matchesPattern(s string, pattern string) (bool, error) {
	// Use filepath.Match which supports shell-style patterns
	return filepath.Match(pattern, s)
}

func getFsFsInDirectory(dir string) ([]string, []string, error) {
	// Ensure dir is clean
	dir = filepath.Clean(dir)

	var folders []string
	var files []string

	children, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}

	for _, file := range children {
		if file.IsDir() {
			folders = append(folders, file.Name())
		} else {
			files = append(files, file.Name())
		}
	}

	return folders, files, nil
}

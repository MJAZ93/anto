package commit

import (
	"anto/model"
	"anto/utils"
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func ValidateCommitMessage(message string, commitFile string) error {
	cm, err := utils.ParseMsk(commitFile)
	if err != nil {
		return errors.New("Error parsing commit.msk, err: " + err.Error())
	}

	if len(message) > cm.MaxLines {
		strLenght := strconv.Itoa(len(message))
		strMaxLines := strconv.Itoa(cm.MaxLines)
		return errors.New("Commit message has " + strLenght + ", when max allowed is " + strMaxLines)
	}

	if len(message) < cm.MinLines {
		strLenght := strconv.Itoa(len(message))
		strMinLines := strconv.Itoa(cm.MinLines)
		return errors.New("Commit message has " + strLenght + ", when min allowed is " + strMinLines)
	}

	matchesAnyRule := false
	for _, regex := range cm.AllowedRules {
		re := regexp.MustCompile(regex)
		if re.MatchString(message) {
			matchesAnyRule = true
		}
	}

	if !matchesAnyRule && len(cm.AllowedRules) > 0 {
		return errors.New("Commit message don't match any allowed rule, commit message instructions: \n\n" + cm.Comment)
	}

	for _, regex := range cm.DisallowedRules {
		re := regexp.MustCompile(regex)
		if re.MatchString(message) {
			return errors.New("Commit message matches: " + regex + " rule, that is not allowed. Commit message instructions: \n\n" + cm.Comment)
		}
	}

	return nil
}

func ValidateCommits() error {
	commitMessage, err := GetCommitMessage(os.Args)
	if err != nil {
		//ClearCommitMessage()
		return errors.New("error getting the commit message. error:" + err.Error())
	}
	// Validate the commit message
	err = ValidateCommitMessage(commitMessage, "commit.msk")
	if err != nil {
		//ClearCommitMessage()
		return errors.New(err.Error())
	}

	return nil
}

// GetCommitMessage reads the commit message from the file path passed as an argument.
func GetCommitMessage(args []string) (string, error) {
	// Ensure the commit message file path is provided
	if len(args) < 2 {
		return "", fmt.Errorf("no commit message file path provided")
	}

	gitFolder, err := utils.GetGitRoot()
	if err != nil {
		return "", nil
	}
	// The commit message file path is the second element in os.Args (args[1])
	commitMsgPath := gitFolder + "/.git/COMMIT_EDITMSG" //args[1]

	// Read the commit message file
	commitMsgBytes, err := os.ReadFile(commitMsgPath)
	if err != nil {
		return "", fmt.Errorf("failed to read commit message: %w", err)
	}

	// Convert the commit message to a string and return
	return string(commitMsgBytes), nil
}

// Parses the .vsk file into the FolderRule structure
func ParseVSKFile(filename string) (model.FolderRule, error) {
	file, err := os.Open(filename)
	if err != nil {
		return model.FolderRule{}, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var root model.FolderRule
	root.Pattern = "root"
	root.MinCount = 1
	root.MaxCount = 1

	var folderStack []*model.FolderRule // Stack of pointers to FolderRule
	var indentStack []int               // Stack of indentation levels

	folderStack = append(folderStack, &root)
	indentStack = append(indentStack, -1) // Root has indentation level -1

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Determine indentation level
		indent := getIndentLevel(line)

		// Trim leading whitespace
		line = strings.TrimLeftFunc(line, unicode.IsSpace)

		// Determine if it's a folder or file
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]") {
			// Folder
			folderRule := model.FolderRule{}
			pattern, minCount, maxCount, err := parsePatternAndCounts(line, '[', ']')
			if err != nil {
				return model.FolderRule{}, err
			}
			folderRule.Pattern = pattern
			folderRule.MinCount = minCount
			folderRule.MaxCount = maxCount

			// Adjust the folder stack based on indentation
			for len(indentStack) > 0 && indentStack[len(indentStack)-1] >= indent {
				indentStack = indentStack[:len(indentStack)-1]
				folderStack = folderStack[:len(folderStack)-1]
			}

			// Add folder to the current folder in the stack
			currentFolder := folderStack[len(folderStack)-1]
			currentFolder.Folders = append(currentFolder.Folders, folderRule)
			// Add the new folder to the stack
			folderStack = append(folderStack, &currentFolder.Folders[len(currentFolder.Folders)-1])
			indentStack = append(indentStack, indent)

		} else if strings.HasPrefix(line, "{") && strings.Contains(line, "}") {
			// File
			fileRule := model.FileRule{}
			pattern, minCount, maxCount, err := parsePatternAndCounts(line, '{', '}')
			if err != nil {
				return model.FolderRule{}, err
			}
			fileRule.Pattern = pattern
			fileRule.MinCount = minCount
			fileRule.MaxCount = maxCount

			// Adjust the folder stack based on indentation
			for len(indentStack) > 0 && indentStack[len(indentStack)-1] >= indent {
				indentStack = indentStack[:len(indentStack)-1]
				folderStack = folderStack[:len(folderStack)-1]
			}
			currentFolder := folderStack[len(folderStack)-1]
			currentFolder.Files = append(currentFolder.Files, fileRule)
		} else {
			fmt.Println("Unknown line format:", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return model.FolderRule{}, fmt.Errorf("error reading file: %v", err)
	}

	return root, nil
}

// Helper functions
func getIndentLevel(line string) int {
	// Count leading spaces
	count := 0
	for _, ch := range line {
		if ch == ' ' {
			count++
		} else {
			break
		}
	}
	indentLevel := count / 4 // Assuming 4 spaces per indent level
	return indentLevel
}

func parsePatternAndCounts(line string, openChar, closeChar rune) (pattern string, minCount, maxCount int, err error) {
	minCount = 0
	maxCount = math.MaxInt32

	// Find the content inside the brackets/braces
	start := strings.IndexRune(line, openChar)
	end := strings.IndexRune(line, closeChar)
	if start == -1 || end == -1 || end <= start {
		pattern = strings.TrimSpace(line)
		return
	}
	pattern = line[start+1 : end]

	// Check for counts after the closing bracket
	rest := strings.TrimSpace(line[end+1:])
	if rest != "" {
		switch {
		case strings.HasPrefix(rest, "="):
			countStr := strings.TrimSpace(rest[1:])
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return pattern, minCount, maxCount, fmt.Errorf("invalid count value: %s", countStr)
			}
			minCount = count
			maxCount = count

		case strings.HasPrefix(rest, ">="):
			countStr := strings.TrimSpace(rest[2:])
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return pattern, minCount, maxCount, fmt.Errorf("invalid count value: %s", countStr)
			}
			minCount = count

		case strings.HasPrefix(rest, ">"):
			countStr := strings.TrimSpace(rest[1:])
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return pattern, minCount, maxCount, fmt.Errorf("invalid count value: %s", countStr)
			}
			minCount = count + 1

		case strings.HasPrefix(rest, "<="):
			countStr := strings.TrimSpace(rest[2:])
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return pattern, minCount, maxCount, fmt.Errorf("invalid count value: %s", countStr)
			}
			maxCount = count

		case strings.HasPrefix(rest, "<"):
			countStr := strings.TrimSpace(rest[1:])
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return pattern, minCount, maxCount, fmt.Errorf("invalid count value: %s", countStr)
			}
			maxCount = count - 1

		case strings.Contains(rest, ".."):
			parts := strings.Split(rest, "..")
			if len(parts) != 2 {
				return pattern, minCount, maxCount, fmt.Errorf("invalid range format: %s", rest)
			}
			minStr := strings.TrimSpace(parts[0])
			maxStr := strings.TrimSpace(parts[1])
			min, err1 := strconv.Atoi(minStr)
			max, err2 := strconv.Atoi(maxStr)
			if err1 != nil || err2 != nil {
				return pattern, minCount, maxCount, fmt.Errorf("invalid range values: %s", rest)
			}
			minCount = min
			maxCount = max

		default:
			// Attempt to parse as a single number (exact count)
			count, err := strconv.Atoi(rest)
			if err == nil {
				minCount = count
				maxCount = count
			} else {
				return pattern, minCount, maxCount, fmt.Errorf("invalid count expression: %s", rest)
			}
		}
	}

	return pattern, minCount, maxCount, nil
}

package parser

import (
	"anto/model"
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func ParseVSKFile(filename string) (model.FolderRule, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return model.FolderRule{}, errors.New("error opening file " + err.Error())
	}
	defer file.Close()

	// Read the file line by line
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
			pattern, minCount, maxCount := parsePatternAndCounts(line, '[', ']')
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
			pattern, minCount, maxCount := parsePatternAndCounts(line, '{', '}')
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
		return model.FolderRule{}, errors.New("error reading file " + err.Error())
	}

	// Output the root FolderRule as JSON
	// jsonOutput, err := json.MarshalIndent(root, "", "    ")
	// if err != nil {
	//	fmt.Println("Error marshalling to JSON:", err)
	//	return model.FolderRule{}, errors.New("Error marshalling file " + err.Error())
	// }
	// fmt.Println(string(jsonOutput))

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

func parsePatternAndCounts(line string, openChar, closeChar rune) (pattern string, minCount, maxCount int) {
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
		// Check for exact count (e.g., "=5")
		if strings.HasPrefix(rest, "=") {
			countStr := rest[1:]
			count, err := strconv.Atoi(countStr)
			if err == nil {
				minCount = count
				maxCount = count
			}
		} else if strings.HasPrefix(rest, ">") { // Minimum count (e.g., ">5")
			countStr := rest[1:]
			count, err := strconv.Atoi(countStr)
			if err == nil {
				minCount = count + 1
				maxCount = math.MaxInt32
			}
		} else if strings.HasPrefix(rest, "<") { // Maximum count (e.g., "<10")
			countStr := rest[1:]
			count, err := strconv.Atoi(countStr)
			if err == nil {
				minCount = 0
				maxCount = count - 1
			}
		} else if strings.Contains(rest, "..") { // Range (e.g., "2..10")
			countRange := strings.Split(rest, "..")
			if len(countRange) == 2 {
				minStr := strings.TrimSpace(countRange[0])
				maxStr := strings.TrimSpace(countRange[1])

				min, err1 := strconv.Atoi(minStr)
				max, err2 := strconv.Atoi(maxStr)

				if err1 == nil {
					minCount = min
				}
				if err2 == nil {
					maxCount = max
				}
			}
		}
	}

	return
}

package generator

import (
	"anto/model"
	"anto/parser"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GenerateFolderAndFileStruct(filename string, basePath string) error {
	// Parse the .vsk file
	rootFolder, err := parser.ParseVSKFile(filename)
	if err != nil {
		return errors.New("Error parsing .vsk file: " + err.Error())
	}

	// Read default content from file-example.msk
	defaultContent, err := os.ReadFile("file-model.msk")
	if err != nil {
		return errors.New("Error reading file-model.msk file: " + err.Error())
	}

	// Ensure basePath exists
	err = os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return errors.New("Error creating base path: " + err.Error())
	}

	err = createStructure(rootFolder, basePath, string(defaultContent))
	if err != nil {
		return errors.New("Error creating structure: " + err.Error())
	}

	fmt.Println("Folder and file structure created successfully under", basePath)
	return nil
}

// Recursively creates the folder and file structure based on the FolderRule
func createStructure(folder model.FolderRule, basePath string, defaultContent string) error {
	// Sanitize the folder name
	folderName := SanitizeName(folder.Pattern)

	// Determine the current path
	currentPath := basePath
	if folderName != "root" {
		currentPath = filepath.Join(basePath, folderName)
	}

	// Ensure the current path exists
	err := os.MkdirAll(currentPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating folder %s: %v", currentPath, err)
	}

	// Create files
	for _, fileRule := range folder.Files {
		// Sanitize the file name
		fileName := SanitizeName(fileRule.Pattern)
		// Handle wildcard patterns by creating a sample file
		if strings.Contains(fileName, "*") || strings.Contains(fileName, "?") {
			fileName = strings.ReplaceAll(fileName, "*", "sample")
			fileName = strings.ReplaceAll(fileName, "?", "x")
		}
		// Add .msk extension
		if !strings.HasSuffix(fileName, ".msk") {
			fileName += ".msk"
		}

		filePath := filepath.Join(currentPath, fileName)
		err := createFile(filePath, defaultContent)
		if err != nil {
			return fmt.Errorf("error creating file %s: %v", filePath, err)
		}
	}

	// Recursively create subfolders
	for _, subFolder := range folder.Folders {
		err := createStructure(subFolder, currentPath, defaultContent)
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function to create a file with specified content
func createFile(filePath string, content string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// Helper function to sanitize file and folder names
func SanitizeName(name string) string {
	// Replace invalid characters with '_'
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	sanitized := invalidChars.ReplaceAllString(name, "_")
	// Replace leading dots with underscores
	sanitized = strings.TrimLeftFunc(sanitized, func(r rune) bool {
		return r == '.'
	})
	if sanitized == "" {
		sanitized = "_"
	}
	// Remove leading and trailing whitespace
	sanitized = strings.TrimSpace(sanitized)
	return sanitized
}

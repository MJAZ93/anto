package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func GenerateStructure(rootDir string) error {
	// Open the output file
	outFile, err := os.Create("structure.vsk")
	if err != nil {
		return errors.New("error creating output file:" + err.Error())
	}
	defer outFile.Close()

	// Build the structure
	structure, err := buildStructure(rootDir)
	if err != nil {
		return err
	}

	// Write the structure to the .vsk file
	writeStructure(outFile, structure, 0, true)
	return nil
}

// Node represents a file or folder in the structure
type Node struct {
	Name     string
	IsDir    bool
	Children []*Node
	Files    []string // List of file names in this directory
}

// buildStructure builds the hierarchical structure starting from rootDir
func buildStructure(rootDir string) (*Node, error) {
	root := &Node{
		Name:  rootDir,
		IsDir: true,
	}
	err := walkDir(root)
	return root, err
}

// walkDir recursively walks through the directory and builds the tree
func walkDir(node *Node) error {
	entries, err := os.ReadDir(node.Name)
	if err != nil {
		return errors.New("error reading directory: " + err.Error())
	}

	for _, entry := range entries {
		// Skip hidden files and directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		childPath := filepath.Join(node.Name, entry.Name())
		if entry.IsDir() {
			child := &Node{
				Name:  childPath,
				IsDir: true,
			}
			node.Children = append(node.Children, child)
			walkDir(child)
		} else {
			node.Files = append(node.Files, entry.Name())
		}
	}

	// Optional: Sort the children by name
	sort.Slice(node.Children, func(i, j int) bool {
		return strings.ToLower(node.Children[i].Name) < strings.ToLower(node.Children[j].Name)
	})
	return nil
}

// writeStructure writes the structure to the output file with indentation
func writeStructure(outFile *os.File, node *Node, depth int, skipRoot bool) {
	var indent string
	if skipRoot && depth == 0 {
		// Do not print the root node
		indent = ""
	} else {
		indent = strings.Repeat("    ", depth)
		name := filepath.Base(node.Name)
		fmt.Fprintf(outFile, "%s[%s]\n", indent, name)
	}
	// Write grouped file patterns
	if len(node.Files) > 0 {
		var fileIndent string
		if skipRoot && depth == 0 {
			fileIndent = ""
		} else {
			fileIndent = strings.Repeat("    ", depth+1)
		}
		patterns := groupFiles(node.Files)
		for _, pattern := range patterns {
			fmt.Fprintf(outFile, "%s{%s}\n", fileIndent, pattern)
		}
	}
	for _, child := range node.Children {
		if skipRoot && depth == 0 {
			// Keep depth the same for root's children
			writeStructure(outFile, child, depth, false)
		} else {
			// Increment depth for nested folders
			writeStructure(outFile, child, depth+1, false)
		}
	}
}

// groupFiles groups files with similar names into patterns
func groupFiles(files []string) []string {
	// Map from pattern to list of files matching that pattern
	patternMap := make(map[string][]string)

	// Group files by extension
	extMap := make(map[string][]string)
	for _, file := range files {
		ext := filepath.Ext(file)
		base := strings.TrimSuffix(file, ext)
		extMap[ext] = append(extMap[ext], base)
	}

	// For each extension group, find common patterns
	for ext, bases := range extMap {
		if len(bases) == 1 {
			// Only one file with this extension, add as is
			pattern := bases[0] + ext
			patternMap[pattern] = []string{pattern}
		} else {
			// Find common prefix and suffix among bases
			prefix := longestCommonPrefix(bases)
			suffix := longestCommonSuffix(bases)

			var pattern string
			if prefix != "" && suffix != "" && prefix != suffix {
				// Since the middle varies, we represent it with '*'
				pattern = fmt.Sprintf("%s*%s%s", prefix, suffix, ext)
			} else if prefix != "" {
				// Common prefix
				pattern = fmt.Sprintf("%s*%s", prefix, ext)
			} else if suffix != "" {
				// Common suffix
				pattern = fmt.Sprintf("*%s%s", suffix, ext)
			} else {
				// No common prefix or suffix, use wildcard
				pattern = fmt.Sprintf("*%s", ext)
			}

			patternMap[pattern] = append(patternMap[pattern], bases...)
		}
	}

	// Collect patterns
	patterns := []string{}
	for pattern := range patternMap {
		patterns = append(patterns, pattern)
	}

	// Sort patterns for consistency
	sort.Strings(patterns)

	return patterns
}

// longestCommonPrefix finds the longest common prefix among strings
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			if len(prefix) == 0 {
				return ""
			}
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
}

// longestCommonSuffix finds the longest common suffix among strings
func longestCommonSuffix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	suffix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasSuffix(s, suffix) {
			if len(suffix) == 0 {
				return ""
			}
			suffix = suffix[1:]
		}
	}
	return suffix
}

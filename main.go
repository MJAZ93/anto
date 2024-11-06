package main

import (
	"anto/commit"
	"anto/generator"
	"anto/parser"
	"anto/precommit"
	"anto/utils"
	"anto/validator"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please provide more parameters.")
		return
	}

	app := os.Args[1]

	err := processRequest(app)
	if err != nil {
		fmt.Println("ANTO VALIDATION: ", err.Error())
		os.Exit(2)
	}
}

func processRequest(request string) error {
	switch request {
	case "init":
		return initProject()
	case "create-validation":
		return createValidationFile()
	case "create-structure":
		return createStructure()
	case "add-precommit":
		return addPreCommitHook()
	case "commit":
		return validateAll()
	case "validate-commit":
		return validateCommit()
	case "validate-structure":
		return validateStructure()
	case "validate":
		return validate()
	default:
		return nil
	}
}

func initProject() error {
	err := createValidationFile()
	if err != nil {
		return err
	}
	err = createStructure()
	if err != nil {
		return err
	}
	err = addPreCommitHook()
	if err != nil {
		return err
	}

	return nil
}

func createStructure() error {
	rootFolder, err := utils.GetGitRootFolderName()
	if err != nil {
		return err
	}
	return generator.GenerateFolderAndFileStruct("structure.vsk", "../.anto/"+rootFolder)
}

func createValidationFile() error {
	rootFolder, err := utils.GetProjectFileName()
	if err != nil {
		return err
	}
	return generator.GenerateStructure(rootFolder)
}

func validate() error {
	err := validateStructure()
	if err != nil {
		return err
	}

	err = validateCommit()
	if err != nil {
		return err
	}

	return nil
}

func validateStructure() error {
	projectName, err := utils.GetProjectFileName()
	if err != nil {
		return err
	}

	root, err := parser.ParseVSKFile(projectName + "/.anto/structure.vsk")
	if err != nil {
		return err
	}

	return validator.ValidateStructure(root, projectName)
}

func validateCommit() error {
	return commit.ValidateCommits()
}

func validateAll() error {
	return commit.ValidateCommits()
}

func addPreCommitHook() error {
	if err := precommit.AddPreCommitHook(); err != nil {
		return err
	}
	return nil
}

package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func readTemplates(templatesDirectory string) []Template {
	entries, _ := os.ReadDir(templatesDirectory)

	templates := []Template{}

	for _, entry := range entries {
		isEntryATemplate := strings.Contains(entry.Name(), "template")

		if !isEntryATemplate || entry.IsDir() {
			continue
		}

		filePath := filepath.Join(templatesDirectory, entry.Name())
		bytes, _ := os.ReadFile(filePath)

		templateContent := string(bytes)

		variableId := 1

		templateVariables := []templateVariable{}
		nextVariableIndex := strings.Index(templateContent, variableDeclaration)

		for nextVariableIndex > 0 {
			variableNameStartIndex := nextVariableIndex + 1 + len(variableDeclaration)
			variableNameEndIndex := strings.Index(templateContent[variableNameStartIndex:], ")")
			variableName := templateContent[variableNameStartIndex : variableNameStartIndex+variableNameEndIndex]
			templateVariables = append(templateVariables, templateVariable{
				id:    variableId,
				name:  variableName,
				index: nextVariableIndex,
			})
			variableId++

			variableDeclarationLength := len(variableDeclaration) + 2 + len(variableName)
			newIndex := strings.Index(templateContent[nextVariableIndex+variableDeclarationLength:], variableDeclaration)
			if newIndex < 0 {
				nextVariableIndex = -1
			} else {
				nextVariableIndex += variableDeclarationLength + newIndex
			}
		}

		template := Template{
			name:      entry.Name(),
			template:  string(bytes),
			variables: templateVariables,
		}

		templates = append(templates, template)
	}

	return templates
}

func readVariables(variableFileName string) map[string]string {
	bytes, _ := os.ReadFile(variableFileName)

	variablesFileContent := string(bytes)

	variables := map[string]string{}

	for index, line := range strings.Split(variablesFileContent, "\n") {

		trimedLine := strings.TrimSpace(line)

		// Skip empty lines
		if len(trimedLine) == 0 {
			continue
		}

		nameAndValue := strings.Split(line, "=")

		if len(nameAndValue) != 2 {
			panic(fmt.Sprintf("Variables file is incorrect at line: %v ", index))
		}

		name := strings.TrimSpace(nameAndValue[0])
		value := strings.TrimSpace(nameAndValue[1])

		variables[name] = value
	}

	return variables
}

func validateVariablesPath(pathToVariables string) []string {
	errors := []string{}

	info, err := os.Stat(pathToVariables)

	if os.IsNotExist(err) {
		errors = append(errors, "Path to variables points to void")
	} else if os.IsPermission(err) {
		errors = append(errors, "Missing permissions to access specified variable file")
	} else if err != nil {
		errors = append(errors, fmt.Sprintf("Uknown error occurred: %v", err))
	}

	if info != nil {
		if info.IsDir() {
			errors = append(errors, "Path to variables must point to file")
		}
	}

	return errors
}

func validateTemplatesPath(pathToTemplates string) []string {
	errors := []string{}

	_, err := os.Stat(pathToTemplates)

	if os.IsNotExist(err) {
		errors = append(errors, "Path to templates points to void")
	} else if os.IsPermission(err) {
		errors = append(errors, "Missing permissions to access specified templates file/directory")
	} else if err != nil {
		errors = append(errors, fmt.Sprintf("Uknown error occurred: %v", err))
	}

	return errors
}

func validateOutputPath(outputPath string) []string {
	errors := []string{}
	info, err := os.Stat(outputPath)

	if os.IsNotExist(err) {
		return errors
	}

	if os.IsPermission(err) {
		errors = append(errors, "Missing permission to access specified output directory")
	} else if err != nil {
		errors = append(errors, fmt.Sprintf("Uknown error occurred: %v", err))
	} else if !info.IsDir() {
		errors = append(errors, "OutputPath points to file")
	}

	return errors
}

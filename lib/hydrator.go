package hydrator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	settingFileName        = "out.settings"
	defaultOutputDirectory = "templates_filled"
	variableDeclaration    = "@hydrate"
)

type templateVariable struct {
	id    int
	name  string
	index int
}

type Template struct {
	name      string
	template  string
	variables []templateVariable
}

func Hydrate(pathToVariables string, pathToTemplates string, outputPath string) {
	homePath := os.Getenv("HOME")

	pathToVariables = strings.ReplaceAll(pathToVariables, "$HOME", homePath)
	pathToTemplates = strings.ReplaceAll(pathToTemplates, "$HOME", homePath)
	outputPath = strings.ReplaceAll(outputPath, "$HOME", homePath)

	outputPathValidationResult := validateOutputPath(outputPath)
	variablesPathValidationResult := validateVariablesPath(pathToVariables)
	templatesPathValidationResult := validateTemplatesPath(pathToTemplates)

	validationErrors := []string{}
	validationErrors = append(validationErrors, outputPathValidationResult...)
	validationErrors = append(validationErrors, variablesPathValidationResult...)
	validationErrors = append(validationErrors, templatesPathValidationResult...)

	if len(validationErrors) != 0 {
		println("There are validation errors:")
		for index, message := range validationErrors {
			fmt.Printf("%v: %v\n", index, message)
		}
		os.Exit(1)
	}

	// outputSettings, err := readOutputSettings(pathToTemplates, homePath)
	// if err != nil {
	// 	println(err)
	// 	os.Exit(1)
	// } else {
	// 	println("Output settings registered")
	// }

	variables := readVariables(pathToVariables)
	templates := readTemplates(pathToTemplates)
	filledTemplates := fillTemplates(variables, &templates)

	fmt.Println(filledTemplates)

	// saveFilledTemplates(&filledTemplates, outputSettings)
}

func readOutputSettings(pathToTemplates string, homePath string) (map[string]string, error) {
	_, err := os.Stat(pathToTemplates)

	settings := map[string]string{}

	if os.IsNotExist(err) {
		return settings, nil
	} else if err != nil {
		return settings, errors.New(fmt.Sprintf("Failed to read template directory %v", err))
	}
	pathToSettings := filepath.Join(pathToTemplates, settingFileName)
	_, err = os.Stat(pathToSettings)

	if os.IsNotExist(err) {
		return settings, nil
	} else if err != nil {
		return settings, errors.New(fmt.Sprintf("Failed to read settings file %v", err))
	}

	settingsFileBytes, err := os.ReadFile(pathToSettings)
	if err != nil {
		return settings, errors.New(fmt.Sprintf("Failed to read settings file %v", err))
	}

	settingsFileContent := string(settingsFileBytes)

	settingsFileLines := strings.Split(settingsFileContent, "\n")

	for _, line := range settingsFileLines {
		optionData := strings.Split(line, "=")

		if len(optionData) != 2 {
			continue
		}

		name := strings.TrimSpace(optionData[0])
		value := strings.TrimSpace(optionData[1])
		value = strings.ReplaceAll(value, "$HOME", homePath)

		parentDirectory := filepath.Dir(value)

		info, err := os.Stat(parentDirectory)

		if os.IsNotExist(err) {
			fmt.Printf("Parent directory for path %v doesn't exists\n", value)
			continue
		} else if err != nil {
			fmt.Printf("Unexpected error %v\n", value)
			continue
		} else if !info.IsDir() {
			fmt.Printf("Parent is not a directory for %v\n", value)
			continue
		}

		settings[name] = value
	}

	return settings, nil
}

func saveFilledTemplates(templates *[]Template, settings map[string]string) {
	for _, templateEntry := range *templates {

		outputPath := settings[templateEntry.name]
		// If user did not specify path for current file then generate default output path
		if outputPath == "" {
			templateSuffixIndex := strings.Index(templateEntry.name, ".template")
			outputPath = filepath.Join(defaultOutputDirectory, templateEntry.name[0:templateSuffixIndex])
		}
		parentDirectory := filepath.Dir(outputPath)
		info, err := os.Stat(parentDirectory)

		if os.IsNotExist(err) {
			os.Mkdir(parentDirectory, 0744)
		} else if !info.IsDir() {
			println("filled_templates is a file")
			os.Exit(1)
		}

		err = os.WriteFile(outputPath, []byte(templateEntry.template), 0744)
		if err != nil {
			fmt.Printf("%v", err)
		}
	}
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

func fillTemplates(variables map[string]string, templates *[]Template) []Template {
	filledTemplates := []Template{}
	for _, templateEntry := range *templates {

		templateContent := templateEntry.template

		for i := len(templateEntry.variables) - 1; i >= 0; i-- {
			entryVariable := templateEntry.variables[i]

			variableValue := variables[entryVariable.name]

			if variableValue == "" {
				fmt.Printf("Template: '%v' contains variable: '%v' that is missing from input\n", templateEntry.name, entryVariable.name)
				os.Exit(1)
			}

			variableDeclarationLength := len(variableDeclaration) + 2 + len(entryVariable.name)
			templateContent = templateContent[:entryVariable.index] + variableValue + templateContent[entryVariable.index+variableDeclarationLength:]

		}

		filledTemplates = append(filledTemplates, Template{
			name:     templateEntry.name,
			template: templateContent,
		})

	}

	return filledTemplates
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

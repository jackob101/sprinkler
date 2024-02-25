package lib

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

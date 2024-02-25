package hydrator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	settingFileName        = "sprinkler.settings"
	defaultOutputDirectory = "templates_filled"
)

type Template struct {
	name     string
	template string
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

	outputSettings, err := readOutputSettings(pathToTemplates)
	if err != nil {
		println(err)
		os.Exit(1)
	} else {
		println("Output settings registered")
	}

	variables := readVariables(pathToVariables)
	templates := readTemplates(pathToTemplates)
	filledTemplates := fillTemplates(variables, &templates)

	saveFilledTemplates(&filledTemplates, outputSettings)
}

func readOutputSettings(pathToTemplates string) (map[string]string, error) {
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

		template := Template{
			name:     entry.Name(),
			template: string(bytes),
		}

		templates = append(templates, template)
	}

	return templates
}

func fillTemplates(variables map[string]string, templates *[]Template) []Template {
	filledTemplates := []Template{}
	for _, templateEntry := range *templates {

		templateContent := templateEntry.template

		for key, variable := range variables {
			println("replacing ", key, variable)
			templateContent = strings.ReplaceAll(templateContent, "$"+key+"$", variable)
			println(templateContent)
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

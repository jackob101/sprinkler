package lib

import (
	"fmt"
	"os"
	"strings"
)

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func Describe(pathToVariables string, pathToTemplates string) {
	homePath := os.Getenv("HOME")

	pathToVariables = strings.ReplaceAll(pathToVariables, "$HOME", homePath)
	pathToTemplates = strings.ReplaceAll(pathToTemplates, "$HOME", homePath)

	variablesPathValidationResult := validateVariablesPath(pathToVariables)
	templatesPathValidationResult := validateTemplatesPath(pathToTemplates)

	validationErrors := []string{}
	validationErrors = append(validationErrors, variablesPathValidationResult...)
	validationErrors = append(validationErrors, templatesPathValidationResult...)

	if len(validationErrors) != 0 {
		println("There are validation errors:")
		for index, message := range validationErrors {
			fmt.Printf("%v: %v\n", index, message)
		}
		os.Exit(1)
	}

	variables := readVariables(pathToVariables)
	templates := readTemplates(pathToTemplates)

	for _, templateEntry := range templates {
		fmt.Printf("Template: %v\n", templateEntry.name)

		if len(templateEntry.variables) == 0 {
			fmt.Printf("%36s \n\n", "There are no declared variables")
			continue
		}

		uniqueVariables := []string{}

		for _, variable := range templateEntry.variables {
			isUnique := true
			for _, uniqueVariable := range uniqueVariables {
				if uniqueVariable == variable.name {
					isUnique = false
					break
				}
			}
			if isUnique {
				uniqueVariables = append(uniqueVariables, variable.name)
			}
		}

		for _, variable := range uniqueVariables {

			variableValue := variables[variable]

			if variableValue == "" {
				fmt.Printf("%-10vErr%v %v\n", colorRed, colorReset, variable)
			} else {
				fmt.Printf("%-10vOk%v %v -> %v\n", colorGreen, colorReset, variable, variableValue)
			}
		}

		println()
	}
}

package hydrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Template struct {
	name     string
	template string
}

func Hydrate() {
	// templatesDirectory := "templates"
	variablesFile := "input.variables"
	variables := readVariables(variablesFile)

	println("Starting hydration")

	fmt.Println(variables)

	templates := readTemplates("templates")

	fmt.Println(templates)

	filledTemplates := fillTemplates(variables, &templates)

	fmt.Println(filledTemplates)
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

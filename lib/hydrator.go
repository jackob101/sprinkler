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

func Hydrate(pathToVariables string, pathToTemplates string) {
	validateVariablesPath(pathToVariables)
	validateTemplatesPath(pathToTemplates)

	variables := readVariables(pathToVariables)
	templates := readTemplates(pathToTemplates)
	filledTemplates := fillTemplates(variables, &templates)

	saveFilledTemplates(&filledTemplates)
}

func saveFilledTemplates(templates *[]Template) {
	pathToSavedTemplates := "filled_templates"

	info, err := os.Stat(pathToSavedTemplates)

	if os.IsNotExist(err) {
		os.Mkdir(pathToSavedTemplates, 0744)
	} else if !info.IsDir() {
		println("filled_templates is a file")
		os.Exit(1)
	}

	for _, templateEntry := range *templates {
		templateSuffixIndex := strings.Index(templateEntry.name, ".template")
		fileName := filepath.Join(pathToSavedTemplates, templateEntry.name[0:templateSuffixIndex])

		err := os.WriteFile(fileName, []byte(templateEntry.template), 0744)
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

func validateVariablesPath(pathToVariables string) {
	path, err := filepath.Abs(pathToVariables)
	if err != nil {
		fmt.Printf("Path to variables failed to convert to absolute. %v \n", err)
		os.Exit(1)
	}

	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		println("Path to variables points to void")
		os.Exit(1)
	} else if err != nil {
		fmt.Println("Uknown error occurred: ", err)
		os.Exit(1)
	}

	if info.IsDir() {
		fmt.Println("Path to variables must point to file not to directory")
		os.Exit(1)
	}
}

func validateTemplatesPath(pathToTemplates string) {
	path, err := filepath.Abs(pathToTemplates)
	if err != nil {
		fmt.Printf("Path to templates failed to convert to absolute. %v \n", err)
		os.Exit(1)
	}

	_, err = os.Stat(path)

	if os.IsNotExist(err) {
		println("Path to templates points to void")
		os.Exit(1)
	} else if err != nil {
		fmt.Println("Uknown error occurred: ", err)
		os.Exit(1)
	}
}

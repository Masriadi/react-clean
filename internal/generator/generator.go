package generator

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masriadi/react-clean/internal/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Embed template files
//
//go:embed templates/repository.tmpl
var repositoryTemplate string

//go:embed templates/entity.tmpl
var entityTemplate string

//go:embed templates/page.tmpl
var pageTemplate string

//go:embed templates/component.tmpl
var componentTemplate string

//go:embed templates/redux_slice.tmpl
var reduxSliceTemplate string

//go:embed templates/redux_thunk.tmpl
var reduxThunkTemplate string

//go:embed templates/routes.tmpl
var routesTemplate string

// TemplateInfo holds the template and its output path
type TemplateInfo struct {
	OutputPath string
	Content    string
}

// Define directory constants
const (
	repositoriesDir = "src/data/repositories"
	entitiesDir     = "src/domain/entities"
	pagesDir        = "src/presentation/pages"
	componentsDir   = "src/presentation/components"
	reduxDir        = "src/presentation/redux"
	routesDir       = "src/presentation/routes"
)

// GenerateStructure generates the directory structure and files for the given entity name
func GenerateStructure(entityName string, moduleName string) error {
	if entityName == "" || moduleName == "" {
		return fmt.Errorf("entityName and moduleName cannot be empty")
	}

	// Show log message
	utils.LogInfo(fmt.Sprintf("Generating structure for entity: %s", entityName))

	// Convert entity name to title case
	entityName = utils.StringToEntityName(entityName)

	// Convert entity name to lowercase
	instanceEntityName := utils.StringToInstanceName(entityName)

	directories := []string{
		repositoriesDir,
		entitiesDir,
		pagesDir,
		reduxDir,
		routesDir,
	}

	// Create directories
	if err := createDirectories(directories); err != nil {
		return err
	}

	// Define templates
	templates := []TemplateInfo{
		{fmt.Sprintf("%s/%s/%sRepository.ts", repositoriesDir, instanceEntityName, entityName), repositoryTemplate},
		{fmt.Sprintf("%s/%s/%s.ts", entitiesDir, instanceEntityName, entityName), entityTemplate},
		{fmt.Sprintf("%s/%s/%sPage.tsx", pagesDir, instanceEntityName, entityName), pageTemplate},
		{fmt.Sprintf("%s/%s/Edit%sModal.tsx", componentsDir, instanceEntityName, entityName), componentTemplate},
		{fmt.Sprintf("%s/%s/%sSlice.ts", reduxDir, instanceEntityName, instanceEntityName), reduxSliceTemplate},
		{fmt.Sprintf("%s/%s/%sThunk.ts", reduxDir, instanceEntityName, instanceEntityName), reduxThunkTemplate},
		{fmt.Sprintf("%s/%s/%sRouter.tsx", routesDir, instanceEntityName, entityName), routesTemplate},
	}

	// Generate files from templates
	for _, tmplInfo := range templates {
		if err := generateFile(
			tmplInfo.OutputPath,
			tmplInfo.Content,
			instanceEntityName,
			entityName,
			moduleName,
		); err != nil {
			return err
		}
	}

	utils.LogSuccess(fmt.Sprintf("Structure for '%s' generated successfully.", instanceEntityName))
	return nil
}

// createDirectories creates the specified directories
func createDirectories(directories []string) error {
	for _, dir := range directories {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}
	return nil
}

// generateFile generates a file from the given template content and entity name
func generateFile(
	outputPathTemplate,
	templateContent,
	instanceEntityName string,
	entityName string,
	moduleName string,
) error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": cases.Title(language.English).String,
	}

	// Process output path template
	fileNameTmpl, err := template.New("filename").Funcs(funcMap).Parse(outputPathTemplate)
	if err != nil {
		return fmt.Errorf("error parsing filename template: %w", err)
	}

	var fileNameBuf strings.Builder
	if err := fileNameTmpl.Execute(&fileNameBuf, map[string]string{"Entity": instanceEntityName}); err != nil {
		return fmt.Errorf("error generating filename: %w", err)
	}
	outPath := fileNameBuf.String()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %w", filepath.Dir(outPath), err)
	}

	// Check if the file already exists
	if _, err := os.Stat(outPath); err == nil {
		if !askOverwrite(outPath) {
			fmt.Printf("Skipping file: %s\n", outPath)
			return nil
		}
	}

	// Parse and execute template
	tmpl, err := template.New("template").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", outPath, err)
	}
	defer outFile.Close()

	data := map[string]string{
		"LowerEntity": instanceEntityName,
		"Entity":      entityName,
		"Module":      moduleName,
	}
	if err := tmpl.Execute(outFile, data); err != nil {
		return fmt.Errorf("error writing to file %s: %w", outPath, err)
	}

	return nil
}

// askOverwrite prompts the user to decide whether to overwrite an existing file
func askOverwrite(filePath string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("File %s already exists. Do you want to overwrite it? (y/n): ", filePath)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)
	return strings.ToLower(response) == "y"
}

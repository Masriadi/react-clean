package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Masriadi/react-clean/internal/generator"
	"github.com/Masriadi/react-clean/internal/remover"
	"github.com/Masriadi/react-clean/internal/utils"
)

// main function to execute commands
func main() {
	if len(os.Args) < 3 {
		utils.LogInfo("Usage: react-clean <command> <name>")
		return
	}

	command := os.Args[1]
	name := os.Args[2]

	// Get package name from package.json
	packageName, err := getPackageNameFromPackageJSON()
	if err != nil {
		utils.LogError(fmt.Sprintf("Error reading package.json: %v", err))
		return
	}
	utils.LogInfo(fmt.Sprintf("Package Name: %s", packageName))

	switch command {
	case "gen":
		if err := generator.GenerateStructure(name, packageName); err != nil {
			utils.LogError(fmt.Sprintf("Error generating structure: %v", err))
		} else {
			utils.LogSuccess(fmt.Sprintf("Structure for '%s' generated successfully.", name))
		}
	case "remove":
		if err := remover.RemoveStructure(name); err != nil {
			utils.LogError(fmt.Sprintf("Error removing structure: %v", err))
		} else {
			utils.LogSuccess(fmt.Sprintf("Structure for '%s' removed successfully.", name))
		}
	default:
		utils.LogError("Unknown command. Use 'gen' or 'remove'.")
	}
}

// getPackageNameFromPackageJSON reads the package.json and returns the "name" field
func getPackageNameFromPackageJSON() (string, error) {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return "", err
	}

	var pkg struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", err
	}

	return pkg.Name, nil
}

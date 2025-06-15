package remover

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masriadi/react-clean/internal/utils"
)

// RemoveStructure removes the generated directories for a given entity
func RemoveStructure(entityName string) error {
	if entityName == "" {
		return fmt.Errorf("entity name cannot be empty")
	}

	// Normalisasi nama
	entityName = utils.StringToEntityName(entityName)
	dirName := utils.StringToDirName(entityName)

	// Daftar direktori yang akan dihapus
	dirs := []string{
		filepath.Join("src/data/repositories", dirName),
		filepath.Join("src/domain/entities", dirName),
		filepath.Join("src/presentation/pages", dirName),
		filepath.Join("src/presentation/components", dirName),
		filepath.Join("src/presentation/stores", dirName),
		filepath.Join("src/presentation/routes", dirName),
	}

	utils.LogInfo("The following directories will be removed:")
	for _, dir := range dirs {
		utils.LogInfo(dir)
	}

	// Konfirmasi pengguna
	fmt.Print("Are you sure you want to proceed? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)
	if strings.ToLower(response) != "y" {
		utils.LogInfo("Removal canceled.")
		return nil
	}

	// Hapus direktori
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("error removing directory %s: %w", dir, err)
		}
	}

	utils.LogSuccess(fmt.Sprintf("Directories for '%s' removed successfully.", entityName))
	return nil
}

// Package requirements provides checks for project requirements documents.
package requirements

import (
	"fmt"
	"os"
	"path/filepath"
)

const requirementsDirectory = "doc/requirements"

// HasDocument checks if a requirements document exists in the project.
func HasDocument(projectRoot string) (bool, error) {
	entries, err := os.ReadDir(filepath.Join(projectRoot, requirementsDirectory))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("read requirements directory: %w", err)
	}

	for _, entry := range entries {
		if entry.Type().IsRegular() && filepath.Ext(entry.Name()) == ".md" {
			return true, nil
		}
	}
	return false, nil
}

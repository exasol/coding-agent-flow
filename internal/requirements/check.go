// Package requirements provides checks for project requirements documents.
package requirements

import (
	"os"
	"path/filepath"

	exaerror "github.com/exasol/error-reporting-go"
)

const requirementsDirectory = "doc/requirements"

// HasDocument checks if a requirements document exists in the project.
func HasDocument(projectRoot string) (bool, error) {
	// [impl~requirements-document-detection~1 -> req~requirements-document-detection~1]
	entries, err := os.ReadDir(filepath.Join(projectRoot, requirementsDirectory))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, exaerror.New("E-CAF-1").
			Message("Failed to read requirements directory: {{error|uq}}.").
			Parameter("error", err).
			Mitigation("Ensure that doc/requirements is a readable directory.")
	}

	for _, entry := range entries {
		if entry.Type().IsRegular() && isSupportedExtension(filepath.Ext(entry.Name())) {
			return true, nil
		}
	}
	return false, nil
}

func isSupportedExtension(extension string) bool {
	switch extension {
	case ".md", ".rst":
		return true
	default:
		return false
	}
}

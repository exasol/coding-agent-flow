package codevalidation_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const codeValidationPackage = "github.com/exasol/coding-agent-flow/internal/codevalidation/"

// forbiddenImports maps each layer to the import path prefixes it avoids.
var forbiddenImports = map[string][]string{
	"domain":              {"os", "net", "github.com/", codeValidationPackage},
	"application/port":    {"os", "net", codeValidationPackage + "application/service", codeValidationPackage + "adapter"},
	"application/service": {"os", "net", codeValidationPackage + "adapter"},
	"adapter/lychee":      {codeValidationPackage + "application/service"},
}

func TestLayersDependInward(t *testing.T) {
	for layer, forbidden := range forbiddenImports {
		files, err := filepath.Glob(filepath.Join(layer, "*.go"))
		if err != nil || len(files) == 0 {
			t.Fatalf("layer %s has no Go files: %v", layer, err)
		}
		for _, file := range files {
			if strings.HasSuffix(file, "_test.go") {
				continue
			}
			violations, err := forbiddenImportsIn(file, forbidden)
			if err != nil {
				t.Fatal(err)
			}
			for _, path := range violations {
				t.Errorf("%s imports %s, which layer %s avoids", file, path, layer)
			}
		}
	}
}

func forbiddenImportsIn(file string, forbidden []string) ([]string, error) {
	source, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	syntax, err := parser.ParseFile(token.NewFileSet(), file, source, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	var violations []string
	for _, spec := range syntax.Imports {
		path := strings.Trim(spec.Path.Value, `"`)
		if slices.ContainsFunc(forbidden, func(prefix string) bool { return matchesPrefix(path, prefix) }) {
			violations = append(violations, path)
		}
	}
	return violations, nil
}

// matchesPrefix reports whether path equals prefix or lies below it.
func matchesPrefix(path, prefix string) bool {
	return path == prefix || strings.HasPrefix(path, prefix+"/") || (strings.HasSuffix(prefix, "/") && strings.HasPrefix(path, prefix))
}

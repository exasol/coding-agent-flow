package validationentrypoint_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const validationEntryPointPackage = "github.com/exasol/coding-agent-flow/internal/validationentrypoint/"

// forbiddenImports maps each layer folder to the import path prefixes it avoids.
var forbiddenImports = map[string][]string{
	"domain":              {"os", "net", "github.com/"},
	"application/port":    {"os", "net", validationEntryPointPackage + "application/service", validationEntryPointPackage + "adapter"},
	"application/service": {"os", "net", validationEntryPointPackage + "adapter"},
	"adapter/*":           {validationEntryPointPackage + "application/service"},
}

// TestLayersDependInward checks that each layer avoids the outer layers, and that the domain, the port, and the service avoid os and net.
func TestLayersDependInward(t *testing.T) {
	for pattern, forbidden := range forbiddenImports {
		files, err := filepath.Glob(filepath.Join(pattern, "*.go"))
		if err != nil || len(files) == 0 {
			t.Fatalf("layer %s has no Go files: %v", pattern, err)
		}
		for _, file := range files {
			if strings.HasSuffix(file, "_test.go") {
				continue
			}
			for _, path := range importsOf(t, file) {
				if slices.ContainsFunc(forbidden, func(prefix string) bool { return matches(path, prefix) }) {
					t.Errorf("%s imports %s, which layer %s avoids", file, path, pattern)
				}
			}
		}
	}
}

// internalPackage is the root of the CAF domains.
const internalPackage = "github.com/exasol/coding-agent-flow/internal/"

// TestOtherDomainsThroughServiceAndDomain checks that each layer uses another domain only through its application service and its domain classes. The root package wires the adapters of each domain.
func TestOtherDomainsThroughServiceAndDomain(t *testing.T) {
	err := filepath.WalkDir(".", func(file string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(file, ".go") || strings.HasSuffix(file, "_test.go") {
			return err
		}
		for _, path := range importsOf(t, file) {
			if otherDomainInternals(path) {
				t.Errorf("%s imports %s; use the application service or the domain classes of that domain", file, path)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// otherDomainInternals reports whether path is an adapter or a port of another domain under internal.
func otherDomainInternals(path string) bool {
	rest, inInternal := strings.CutPrefix(path, internalPackage)
	if !inInternal || strings.HasPrefix(path, validationEntryPointPackage) {
		return false
	}
	_, layer, _ := strings.Cut(rest, "/")
	return matches(layer, "adapter") || matches(layer, "application/port")
}

// matches reports whether path equals prefix or lies below it.
func matches(path, prefix string) bool {
	return path == prefix || strings.HasPrefix(path, prefix+"/") || (strings.HasSuffix(prefix, "/") && strings.HasPrefix(path, prefix))
}

// importsOf returns the import paths of a Go file.
func importsOf(t *testing.T, file string) []string {
	t.Helper()
	source, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	syntax, err := parser.ParseFile(token.NewFileSet(), file, source, parser.ImportsOnly)
	if err != nil {
		t.Fatal(err)
	}
	paths := make([]string, 0, len(syntax.Imports))
	for _, spec := range syntax.Imports {
		paths = append(paths, strings.Trim(spec.Path.Value, `"`))
	}
	return paths
}

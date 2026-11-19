package internal_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/kcmvp/archunit"
)

const internalPackage = "github.com/exasol/coding-agent-flow/internal/"

// globalFolder holds the domain classes that all domains share. Each domain layer may refer to it.
const globalFolder = "global"

// layerFolders maps each layer of a domain to its folder pattern. The suffix "..." matches the folder package and its sub-packages.
var layerFolders = map[string]string{
	"Domain":  "domain...",
	"Port":    "application/port...",
	"Service": "application/service...",
	"Adapter": "adapter...",
}

// innerLayerFolders lists the layer folders that avoid os and net, such as os/exec and net/http.
var innerLayerFolders = []string{"domain", "application/port", "application/service"}

// domains returns the folders under internal that contain a domain folder, such as oftvalidation.
func domains(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	var result []string
	for _, entry := range entries {
		if info, err := os.Stat(filepath.Join(entry.Name(), "domain")); err == nil && entry.IsDir() && info.IsDir() {
			result = append(result, entry.Name())
		}
	}
	if len(result) == 0 {
		t.Fatal("internal has no domain")
	}
	return result
}

func layer(domain, name string) string {
	return domain + "/" + name
}

func hasGlobal() bool {
	info, err := os.Stat(globalFolder)
	return err == nil && info.IsDir()
}

// architecture defines the layers of all domains, and the Global layer if internal/global exists. archunit keeps one architecture per test binary.
func architecture(t *testing.T) archunit.Architecture {
	t.Helper()
	var layers []*archunit.Layer
	for _, domain := range domains(t) {
		for name, folder := range layerFolders {
			layers = append(layers, archunit.ArchLayer(layer(domain, name), internalPackage+domain+"/"+folder))
		}
	}
	if hasGlobal() {
		layers = append(layers, archunit.ArchLayer("Global", internalPackage+globalFolder))
	}
	return archunit.ArchUnit(layers...)
}

// TestLayersSelectPackages checks with archunit that each layer of each domain selects at least one package. archunit passes every rule on an empty layer.
func TestLayersSelectPackages(t *testing.T) {
	architecture(t)
	for _, domain := range domains(t) {
		for name := range layerFolders {
			if len(archunit.Layers(layer(domain, name)).Packages().Objects()) == 0 {
				t.Errorf("layer %s selects no packages", layer(domain, name))
			}
		}
	}
}

// TestLayersDependInward checks with archunit that the layers of each domain depend inward: the domain refers only to itself, the Global layer, and the standard library, the service avoids the adapters, and the adapters avoid the service. Go rejects import cycles, which covers the other directions.
func TestLayersDependInward(t *testing.T) {
	arch := architecture(t)
	for _, domain := range domains(t) {
		allowedByDomain := []string{layer(domain, "Domain")}
		if hasGlobal() {
			allowedByDomain = append(allowedByDomain, "Global")
		}
		rules := map[string]archunit.Rule{
			"domain refers only to itself, Global, and the standard library": archunit.Layers(layer(domain, "Domain")).ShouldOnlyRefer(archunit.Layers(allowedByDomain...)),
			"service avoids adapter": archunit.Layers(layer(domain, "Service")).ShouldNotRefer(archunit.Layers(layer(domain, "Adapter"))),
			"adapter avoids service": archunit.Layers(layer(domain, "Adapter")).ShouldNotRefer(archunit.Layers(layer(domain, "Service"))),
		}
		for name, rule := range rules {
			if err := arch.Validate(rule); err != nil {
				t.Errorf("%s: %s: %v", domain, name, err)
			}
		}
	}
}

// TestInnerLayersAvoidInputOutput checks with go/parser that the domain, the port, and the service of each domain avoid os and net, such as os/exec and net/http. archunit applies its rules to module packages only, so this test reads the imports of each file.
func TestInnerLayersAvoidInputOutput(t *testing.T) {
	for _, domain := range domains(t) {
		for _, folder := range innerLayerFolders {
			for _, file := range goFiles(t, filepath.Join(domain, folder)) {
				for _, path := range importsOf(t, file) {
					if slices.ContainsFunc([]string{"os", "net"}, func(prefix string) bool { return path == prefix || strings.HasPrefix(path, prefix+"/") }) {
						t.Errorf("%s imports %s, which the inner layers avoid", file, path)
					}
				}
			}
		}
	}
}

// goFiles returns the Go files of a folder and its sub-folders, apart from test files.
func goFiles(t *testing.T, folder string) []string {
	t.Helper()
	var files []string
	err := filepath.WalkDir(folder, func(path string, entry fs.DirEntry, err error) error {
		if err == nil && !entry.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}
		return err
	})
	if err != nil || len(files) == 0 {
		t.Fatalf("%s has no Go files: %v", folder, err)
	}
	return files
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

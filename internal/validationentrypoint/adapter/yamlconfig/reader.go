// Package yamlconfig implements ConfigurationReader for the .caf.yaml file of a Target Project.
package yamlconfig

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	exaerror "github.com/exasol/error-reporting-go"
	"gopkg.in/yaml.v3"

	"github.com/exasol/coding-agent-flow/internal/global"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/application/port"
)

// FileName is the configuration file in the Target Project root.
const FileName = ".caf.yaml"

// currentVersion is the supported configuration version.
const currentVersion = 1

// configuration is the content of .caf.yaml.
type configuration struct {
	Version    int      `yaml:"version"`
	Validators []string `yaml:"validators"`
}

type Reader struct{}

var _ port.ConfigurationReader = Reader{}

func NewReader() Reader {
	return Reader{}
}

// ValidationSelection reads the selected Validation Names from .caf.yaml in the Target Project root.
func (Reader) ValidationSelection(_ context.Context, targetProject global.TargetProject) ([]string, error) {
	// [impl -> req~validation-selection~1]
	path := filepath.Join(targetProject.Path(), FileName)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, exaerror.New("E-CAF-5").
			Message("Failed to read the CAF configuration {{path}}: {{error|uq}}.").
			Parameter("path", path).
			Parameter("error", err).
			Mitigation("Add a readable " + FileName + " file to the Target Project root.")
	}
	selection, err := parse(content)
	if err != nil {
		return nil, exaerror.New("E-CAF-6").
			Message("Invalid CAF configuration {{path}}: {{error|uq}}.").
			Parameter("path", path).
			Parameter("error", err).
			Mitigation(fmt.Sprintf("Use configuration version %d and select at least one validator with a unique name.", currentVersion))
	}
	return selection, nil
}

// parse reads one YAML document with known fields only, the supported version, and unique, non-empty Validation Names.
func parse(content []byte) ([]string, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	var result configuration
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, errors.New("the file has more than one YAML document")
	}
	if result.Version != currentVersion {
		return nil, fmt.Errorf("version %d is unsupported", result.Version)
	}
	if len(result.Validators) == 0 {
		return nil, errors.New("the file selects no validator")
	}
	seen := make(map[string]bool, len(result.Validators))
	selection := make([]string, 0, len(result.Validators))
	for index, name := range result.Validators {
		name = strings.TrimSpace(name)
		if name == "" {
			return nil, fmt.Errorf("validator %d has no name", index+1)
		}
		if seen[name] {
			return nil, fmt.Errorf("validator %q appears twice", name)
		}
		seen[name] = true
		selection = append(selection, name)
	}
	return selection, nil
}

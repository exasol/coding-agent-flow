// Package config discovers and loads CAF project configuration.
package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	exaerror "github.com/exasol/error-reporting-go"
	"gopkg.in/yaml.v3"
)

const (
	// FileName is the single project-configuration filename recognized by CAF.
	FileName = ".caf.yaml"
	// CurrentVersion is the supported configuration schema version.
	CurrentVersion = 1
)

// Config defines the versioned validation policy for one project.
type Config struct {
	Version    int      `yaml:"version"`
	Validators []string `yaml:"validators"`
}

// Project is a discovered project and its validated CAF configuration.
type Project struct {
	Root       string
	ConfigPath string
	Config     Config
}

// Discover searches from start upwards for FileName. The nearest file wins and
// its containing directory is the project root.
func Discover(start string) (Project, error) {
	// [impl~caf-project-discovery~1 -> req~caf-project-discovery~1]
	directory, err := filepath.Abs(start)
	if err != nil {
		return Project{}, discoveryError(err)
	}
	info, err := os.Stat(directory)
	if err != nil || !info.IsDir() {
		if err == nil {
			err = fmt.Errorf("start path is not a directory")
		}
		return Project{}, discoveryError(err)
	}

	for {
		configPath := filepath.Join(directory, FileName)
		info, statErr := os.Stat(configPath)
		if statErr == nil {
			if !info.Mode().IsRegular() {
				return Project{}, discoveryError(fmt.Errorf("configuration is not a regular file"))
			}
			loaded, loadErr := load(configPath)
			if loadErr != nil {
				return Project{}, loadErr
			}
			return Project{Root: directory, ConfigPath: configPath, Config: loaded}, nil
		}
		if !os.IsNotExist(statErr) {
			return Project{}, discoveryError(statErr)
		}

		parent := filepath.Dir(directory)
		if parent == directory {
			return Project{}, discoveryError(fmt.Errorf("%s not found", FileName))
		}
		directory = parent
	}
}

func load(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, discoveryError(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	var result Config
	if err := decoder.Decode(&result); err != nil {
		return Config{}, invalidConfigError(err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			err = fmt.Errorf("multiple YAML documents are not supported")
		}
		return Config{}, invalidConfigError(err)
	}
	if result.Version != CurrentVersion {
		return Config{}, invalidConfigError(fmt.Errorf("unsupported version %d", result.Version))
	}
	if len(result.Validators) == 0 {
		return Config{}, invalidConfigError(fmt.Errorf("validators must not be empty"))
	}
	seen := make(map[string]struct{}, len(result.Validators))
	for index, validator := range result.Validators {
		validator = strings.TrimSpace(validator)
		if validator == "" {
			return Config{}, invalidConfigError(fmt.Errorf("validator at index %d is empty", index))
		}
		if _, exists := seen[validator]; exists {
			return Config{}, invalidConfigError(fmt.Errorf("validator %q is duplicated", validator))
		}
		seen[validator] = struct{}{}
		result.Validators[index] = validator
	}
	return result, nil
}

func discoveryError(cause error) error {
	return exaerror.New("E-CAF-2").
		Message("Failed to discover CAF project configuration: {{cause|uq}}.").
		Parameter("cause", cause).
		Mitigation("Add a readable .caf.yaml file to the project root.")
}

func invalidConfigError(cause error) error {
	return exaerror.New("E-CAF-3").
		Message("Invalid CAF project configuration: {{cause|uq}}.").
		Parameter("cause", cause).
		Mitigation("Use configuration version 1 and declare at least one unique validator.")
}

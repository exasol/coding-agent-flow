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
	"regexp"
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

// maximumFileSize limits .caf.yaml in bytes.
const maximumFileSize = 64 << 10

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
	content, err := read(path)
	if err != nil {
		return nil, err
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

// read returns the content of .caf.yaml. The file is a regular file of at most maximumFileSize bytes: CAF checks the
// file type before it opens the file, so a symbolic link, a FIFO, or a device stays unopened, and checks the opened file
// again for a replacement in between.
func read(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, readError(path, err)
	}
	if !info.Mode().IsRegular() {
		return nil, notRegularError(path, info.Mode())
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, readError(path, err)
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil {
		return nil, readError(path, err)
	}
	if !opened.Mode().IsRegular() || !os.SameFile(info, opened) {
		return nil, notRegularError(path, opened.Mode())
	}
	content, err := io.ReadAll(io.LimitReader(file, maximumFileSize+1))
	if err != nil {
		return nil, readError(path, err)
	}
	if len(content) > maximumFileSize {
		return nil, exaerror.New("E-CAF-11").
			Message("The CAF configuration {{path}} exceeds {{limit}} bytes.").
			Parameter("path", path).
			Parameter("limit", maximumFileSize).
			Mitigation("Keep " + FileName + " below the limit; it lists Validation Names only.")
	}
	return content, nil
}

func readError(path string, err error) error {
	return exaerror.New("E-CAF-5").
		Message("Failed to read the CAF configuration {{path}}: {{error|uq}}.").
		Parameter("path", path).
		Parameter("error", err).
		Mitigation("Add a readable " + FileName + " file to the Target Project root.")
}

func notRegularError(path string, mode os.FileMode) error {
	return exaerror.New("E-CAF-10").
		Message("The CAF configuration {{path}} is no regular file, its type is {{type}}.").
		Parameter("path", path).
		Parameter("type", mode.Type().String()).
		Mitigation("Make " + FileName + " a regular file in the Target Project root, not a symbolic link.")
}

// parse reads one YAML document with known fields only, the supported version, and unique, non-empty Validation Names.
func parse(content []byte) ([]string, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	var result configuration
	if err := decoder.Decode(&result); err != nil {
		return nil, schemaError(err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, errors.New("the file has more than one YAML document")
	}
	if result.Version != currentVersion {
		return nil, errors.New("the version is unsupported")
	}
	if len(result.Validators) == 0 {
		return nil, errors.New("the file selects no validator")
	}
	positions := make(map[string]int, len(result.Validators))
	selection := make([]string, 0, len(result.Validators))
	for index, name := range result.Validators {
		name = strings.TrimSpace(name)
		if name == "" {
			return nil, fmt.Errorf("validator %d has no name", index+1)
		}
		if first, found := positions[name]; found {
			return nil, fmt.Errorf("validator %d repeats validator %d", index+1, first)
		}
		positions[name] = index + 1
		selection = append(selection, name)
	}
	return selection, nil
}

// parserLine finds the line number at the start of a YAML parser message.
var parserLine = regexp.MustCompile(`^(?:yaml: )?line (\d+):`)

// schemaError describes a YAML decoding error by its kind and line only. The public message contains no content of
// .caf.yaml, because parser messages quote keys and values, and a value can be a credential.
func schemaError(err error) error {
	if errors.Is(err, io.EOF) {
		return errors.New("the file is empty")
	}
	var typeError *yaml.TypeError
	if !errors.As(err, &typeError) || len(typeError.Errors) == 0 {
		return errors.New(atLine(err.Error()) + "the YAML syntax is invalid")
	}
	problem := typeError.Errors[0]
	switch {
	case strings.Contains(problem, "not found in type"):
		return errors.New(atLine(problem) + "the field is unknown")
	case strings.Contains(problem, "already defined"):
		return errors.New(atLine(problem) + "the field appears twice")
	default:
		return errors.New(atLine(problem) + "the value has the wrong type")
	}
}

// atLine returns the "line n: " prefix of a parser message, or an empty string for a message without a line.
func atLine(message string) string {
	if match := parserLine.FindStringSubmatch(message); match != nil {
		return "line " + match[1] + ": "
	}
	return ""
}

package acceptance

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

const atDataFile = "at_data.toml"

var (
	variablePattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)}`)
	valuePattern    = regexp.MustCompile(`\{\{([^{}]+)}}`)
)

// testRun is one [[scenario]] block of at_data.toml.
type testRun struct {
	number      int
	id          string
	description string
	given       Values
	files       []inputFile
	then        Values
}

type inputFile struct {
	Path    string
	Content string
}

// name returns "<id>: <description>", or run-<number> for a run with neither.
func (run testRun) name() string {
	label := strings.Join(slices.DeleteFunc([]string{run.id, run.description}, func(part string) bool { return part == "" }), ": ")
	if label == "" {
		return fmt.Sprintf("run-%d", run.number)
	}
	return label
}

// writeFiles writes the run's input files below directory.
func (run testRun) writeFiles(directory string) error {
	for _, file := range run.files {
		path := filepath.Join(directory, filepath.FromSlash(file.Path))
		if !strings.HasPrefix(path, filepath.Clean(directory)+string(filepath.Separator)) {
			return fmt.Errorf("input file %q leaves the file directory", file.Path)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create input directory: %w", err)
		}
		if err := os.WriteFile(path, []byte(file.Content), 0o600); err != nil {
			return fmt.Errorf("write input file: %w", err)
		}
	}
	return nil
}

type atDataDocument struct {
	Variables map[string]string `toml:"variables"`
	Scenarios []struct {
		ID          string         `toml:"id"`
		Description string         `toml:"description"`
		Given       map[string]any `toml:"given"`
		Then        map[string]any `toml:"then"`
	} `toml:"scenario"`
}

// readAtData reads the runs of at_data.toml in file order.
func readAtData(path string) ([]testRun, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var document atDataDocument
	decoder := toml.NewDecoder(bytes.NewReader(content)).DisallowUnknownFields()
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if len(document.Scenarios) == 0 {
		return nil, fmt.Errorf("%s has no [[scenario]] block", path)
	}
	runs := make([]testRun, len(document.Scenarios))
	for index, scenario := range document.Scenarios {
		if runs[index], err = readRun(index+1, scenario.Given, scenario.Then, document.Variables); err != nil {
			return nil, fmt.Errorf("%s, run %d: %w", path, index+1, err)
		}
		runs[index].id, runs[index].description = scenario.ID, scenario.Description
	}
	return runs, nil
}

// readRun reads one [[scenario]] block and fills its {{name}} and ${name} references.
func readRun(number int, given, then map[string]any, variables map[string]string) (testRun, error) {
	run := testRun{number: number, given: textValues(given), then: textValues(then)}
	delete(run.given, "file")
	for name, value := range run.then {
		filled, err := fill(value, valuePattern, run.given)
		if err != nil {
			return testRun{}, err
		}
		run.then[name] = filled
	}
	files, err := givenFiles(given)
	if err != nil {
		return testRun{}, err
	}
	for _, file := range files {
		if file.Path, err = fillFileText(file.Path, run.given, variables); err != nil {
			return testRun{}, err
		}
		if file.Content, err = fillFileText(file.Content, run.given, variables); err != nil {
			return testRun{}, err
		}
		run.files = append(run.files, file)
	}
	return run, nil
}

func fillFileText(text string, given Values, variables map[string]string) (string, error) {
	text, err := fill(text, variablePattern, variables)
	if err != nil {
		return "", err
	}
	return fill(text, valuePattern, given)
}

func textValues(table map[string]any) Values {
	values := Values{}
	for name, value := range table {
		values[name] = fmt.Sprint(value)
	}
	return values
}

func givenFiles(given map[string]any) ([]inputFile, error) {
	entries, isList := given["file"].([]any)
	if given["file"] != nil && !isList {
		return nil, fmt.Errorf("write file entries as [[scenario.given.file]]")
	}
	files := make([]inputFile, 0, len(entries))
	for _, entry := range entries {
		table, isTable := entry.(map[string]any)
		path, hasPath := table["path"].(string)
		content, _ := table["content"].(string)
		if !isTable || !hasPath || path == "" {
			return nil, fmt.Errorf("a [[scenario.given.file]] entry has no path")
		}
		files = append(files, inputFile{Path: path, Content: content})
	}
	return files, nil
}

func fill(text string, pattern *regexp.Regexp, values map[string]string) (string, error) {
	var problem error
	filled := pattern.ReplaceAllStringFunc(text, func(reference string) string {
		name := strings.TrimSpace(pattern.FindStringSubmatch(reference)[1])
		value, exists := values[name]
		if !exists {
			problem = fmt.Errorf("%s names no value", reference)
		}
		return value
	})
	return filled, problem
}

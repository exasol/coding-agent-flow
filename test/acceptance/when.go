package acceptance

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

// When binds the When step of a use case to its public call.
//
//	Step     the When step, shown in failure messages
//	Call     a Go expression returning (response, error); it can use ctx, files.Directory(), given, and Inputs
//	Imports  the import specs Call needs, for example `api "example.com/module"`
//	Inputs   named variables for Call
//	Decode   reads the JSON response
type When[Response any] struct {
	Step    string
	Call    string
	Imports []string
	Inputs  []Input
	Decode  func([]byte) (Response, error)
}

// Input is a named variable for Call; Imports are the import specs Expression needs.
type Input struct {
	Name       string
	Expression string
	Imports    []string
}

func DecodeJSON[Response any](output []byte) (Response, error) {
	var response Response
	if err := json.Unmarshal(output, &response); err != nil {
		return response, fmt.Errorf("decode response %q: %w", output, err)
	}
	return response, nil
}

// callProgram is a generated Go program that makes the public call. The tests compile independently of
// the public API, so a missing method shows up as a failing test.
type callProgram struct {
	source string
}

func newCallProgram(call string, imports []string, inputs []Input) callProgram {
	imports = append(slices.Clone(imports), `"context"`, `"encoding/json"`, `"fmt"`, `"os"`)
	var declarations strings.Builder
	for _, input := range inputs {
		fmt.Fprintf(&declarations, "\t%s := %s\n\t_ = %s\n", input.Name, input.Expression, input.Name)
		imports = append(imports, input.Imports...)
	}
	slices.Sort(imports)
	return callProgram{source: fmt.Sprintf(programTemplate,
		"\t"+strings.Join(slices.Compact(imports), "\n\t"), declarations.String(), call)}
}

const programTemplate = `package main

import (
%s
)

// inputFiles is the directory with the run's input files.
type inputFiles string

func (files inputFiles) Directory() string {
	return string(files)
}

// Values holds the Given values for the current scenario.
type Values map[string]string

func (values Values) String(name string) string {
	value, exists := values[name]
	if !exists {
		fail("at_data.toml has no Given value %%q", name)
	}
	return value
}

func fail(format string, arguments ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", arguments...)
	os.Exit(1)
}

func main() {
	ctx := context.Background()
	files := inputFiles(os.Args[1])
	given := Values{}
	if err := json.Unmarshal([]byte(os.Args[2]), &given); err != nil {
		fail("decode Given values: %%v", err)
	}
	_, _, _ = ctx, files, given
%s
	response, err := %s
	if err != nil {
		fail("%%v", err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		fail("%%v", err)
	}
}
`

type compiledCallProgram struct {
	binary string
}

var (
	goDirectivePattern     = regexp.MustCompile(`(?m)^go\s+(\S+)\s*$`)
	moduleDirectivePattern = regexp.MustCompile(`(?m)^module\s+(\S+)\s*$`)
)

// build compiles the program in its own module, which uses the repository's module.
func (program callProgram) build(ctx context.Context, directory string) (compiledCallProgram, error) {
	root := repositoryRoot()
	rootModule, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return compiledCallProgram{}, fmt.Errorf("read go.mod: %w", err)
	}
	modulePath := moduleDirectivePattern.FindSubmatch(rootModule)
	goVersion := goDirectivePattern.FindSubmatch(rootModule)
	if modulePath == nil || goVersion == nil {
		return compiledCallProgram{}, errors.New("go.mod has no module or go directive")
	}
	module := fmt.Sprintf("module callprogram\n\ngo %s\n\nrequire %s v0.0.0\n\nreplace %s => %s\n", goVersion[1], modulePath[1], modulePath[1], root)
	files := map[string][]byte{"go.mod": []byte(module), "main.go": []byte(program.source)}
	if sums, err := os.ReadFile(filepath.Join(root, "go.sum")); err == nil {
		files["go.sum"] = sums
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(directory, name), content, 0o600); err != nil {
			return compiledCallProgram{}, fmt.Errorf("write %s: %w", name, err)
		}
	}
	binary := filepath.Join(directory, "call-program")
	if _, err := runCommand(ctx, directory, "go", "build", "-mod=mod", "-buildvcs=false", "-o", binary, "."); err != nil {
		return compiledCallProgram{}, err
	}
	return compiledCallProgram{binary: binary}, nil
}

// call runs the program on the run's file directory and returns the JSON response.
func (program compiledCallProgram) call(ctx context.Context, fileDirectory string, given Values) ([]byte, error) {
	givenJSON, err := json.Marshal(given)
	if err != nil {
		return nil, fmt.Errorf("encode Given values: %w", err)
	}
	return runCommand(ctx, "", program.binary, fileDirectory, string(givenJSON))
}

// runCommand returns the command's stdout. On failure, it returns an error with the command's stderr.
func runCommand(ctx context.Context, directory string, name string, arguments ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, name, arguments...)
	command.Dir = directory
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return nil, errors.New(message)
	}
	return bytes.TrimSpace(output), nil
}

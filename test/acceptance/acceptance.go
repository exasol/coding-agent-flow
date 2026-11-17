// Package acceptance runs the acceptance test of a use case against its public API.
//
// Each test folder holds at_data.toml (input and expected values per run) and at_test.go (the public call,
// checks, and test function).
package acceptance

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type Suite[Response any] struct {
	when   When[Response]
	checks []Check[Response]
}

// Check verifies a response using the current scenario's Given and Then values.
type Check[Response any] func(response Response, given, then Values) error

func NewSuite[Response any](when When[Response], checks ...Check[Response]) Suite[Response] {
	return Suite[Response]{when: when, checks: checks}
}

// Run runs every [[scenario]] block of the at_data.toml next to the calling test.
func (suite Suite[Response]) Run(t *testing.T) {
	t.Helper()
	if err := suite.requireChecks(); err != nil {
		t.Fatal(err)
	}
	functionName, _, _ := strings.Cut(t.Name(), "/")
	testFile, err := sourceFileOf(functionName)
	if err != nil {
		t.Fatal(err)
	}
	runs, err := readAtData(filepath.Join(filepath.Dir(testFile), atDataFile))
	if err != nil {
		t.Fatal(err)
	}
	program, err := newCallProgram(suite.when.Call, suite.when.Imports, suite.when.Inputs).build(t.Context(), t.TempDir())
	if err != nil {
		t.Fatalf("public call %s does not compile:\n%v", suite.when.Call, err)
	}
	for _, run := range runs {
		t.Run(run.name(), func(t *testing.T) {
			t.Parallel()
			suite.run(t, program, run)
		})
	}
}

func (suite Suite[Response]) requireChecks() error {
	if len(suite.checks) == 0 {
		return errors.New("the suite has no check")
	}
	return nil
}

func (suite Suite[Response]) run(t *testing.T, program compiledCallProgram, run testRun) {
	t.Logf("given %v, then %v", run.given, run.then)
	fileDirectory := t.TempDir()
	if err := run.writeFiles(fileDirectory); err != nil {
		t.Fatalf("Given: %v", err)
	}
	output, err := program.call(t.Context(), fileDirectory, run.given)
	if err != nil {
		t.Fatalf("When %s: %v", suite.when.Step, err)
	}
	response, err := suite.when.Decode(output)
	if err != nil {
		t.Fatalf("When %s: %v", suite.when.Step, err)
	}
	for _, check := range suite.checks {
		if err := runCheck(check, response, run.given, run.then); err != nil {
			t.Errorf("Then: %v", err)
		}
	}
}

// runCheck turns a panic about a missing or malformed value into an error.
func runCheck[Response any](check Check[Response], response Response, given, then Values) (err error) {
	defer func() {
		if problem := recover(); problem != nil {
			err = fmt.Errorf("%v", problem)
		}
	}()
	return check(response, given, then)
}

func sourceFileOf(functionName string) (string, error) {
	callers := make([]uintptr, 64)
	frames := runtime.CallersFrames(callers[:runtime.Callers(1, callers)])
	for {
		frame, more := frames.Next()
		if strings.HasSuffix(frame.Function, "."+functionName) {
			return frame.File, nil
		}
		if !more {
			return "", fmt.Errorf("call Run from test function %s", functionName)
		}
	}
}

func repositoryRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

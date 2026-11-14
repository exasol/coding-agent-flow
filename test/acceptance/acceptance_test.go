package acceptance

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
)

type scenarioID string
type targetProjectPath string
type scenarioCatalogPath string
type fixtureCatalogPath string
type goProgramSource string
type programOutput string
type targetProjectInvocation string

type useCaseLocation struct {
	feature     string
	requirement string
	useCase     string
}

func (location useCaseLocation) scenarioCatalog() scenarioCatalogPath {
	return scenarioCatalogPath(filepath.Join(repositoryRoot(), "doc", "requirements", "scenarios", location.feature, location.requirement, location.useCase, "scenarios.feature"))
}

func (location useCaseLocation) fixtureCatalog() fixtureCatalogPath {
	return fixtureCatalogPath(filepath.Join(filepath.Dir(testSourceFile()), "fixtures", location.feature, location.requirement, location.useCase+".fixture"))
}

type publicBoundaryProgram struct {
	source goProgramSource
}

func newTargetProjectProgram(invocation targetProjectInvocation) publicBoundaryProgram {
	return publicBoundaryProgram{source: goProgramSource(fmt.Sprintf(`package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	codingagentflow "github.com/exasol/coding-agent-flow"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "expected target project path")
		os.Exit(2)
	}
	response, err := %s
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
`, invocation))}
}

type responseDecoder[Response any] func(programOutput) (Response, error)
type responseMatcher[Response any] func(Response, Response, targetProjectFixture) error

type oftElementType string
type oftElementID string

const oftRequirementElementType oftElementType = "req"

type oftValidationResponse struct {
	Valid                 bool                              `json:"valid"`
	OFTElementIDs         map[oftElementType][]oftElementID `json:"oftElementIDs"`
	fixtureElementIDTypes []oftElementType
}

func decodeOFTValidation(output programOutput) (oftValidationResponse, error) {
	var response oftValidationResponse
	if err := json.Unmarshal([]byte(output), &response); err != nil {
		return oftValidationResponse{}, fmt.Errorf("read OFT validation response %q: %w", output, err)
	}
	return response, nil
}

func matchOFTValidation(expected, actual oftValidationResponse, fixture targetProjectFixture) error {
	if actual.Valid != expected.Valid {
		return fmt.Errorf("expected valid %t, got %t", expected.Valid, actual.Valid)
	}
	for elementType, expectedIDs := range expected.OFTElementIDs {
		actualIDs := actual.OFTElementIDs[elementType]
		if !sameElementIDs(expectedIDs, actualIDs) {
			return fmt.Errorf("expected %s element IDs %v, got %v", elementType, expectedIDs, actualIDs)
		}
	}
	for _, elementType := range expected.fixtureElementIDTypes {
		expectedIDs := fixture.elementIDs(elementType)
		actualIDs := actual.OFTElementIDs[elementType]
		if !sameElementIDs(expectedIDs, actualIDs) {
			return fmt.Errorf("expected %s element IDs %v, got %v", elementType, expectedIDs, actualIDs)
		}
	}
	return nil
}

var oftElementIDPattern = regexp.MustCompile(`([a-z]+)~([A-Za-z0-9-]+)~([0-9]+)`)

type scenarioBinding[Response any] struct {
	expected Response
}

type targetProjectFile struct {
	path    string
	content string
}

type targetProjectFixture struct {
	files []targetProjectFile
}

func (fixture targetProjectFixture) create(t *testing.T) targetProjectPath {
	t.Helper()
	if len(fixture.files) == 0 {
		t.Fatal("missing target project fixture files")
	}
	projectDirectory := t.TempDir()
	for _, file := range fixture.files {
		filePath := filepath.Join(projectDirectory, file.path)
		if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
			t.Fatalf("create fixture directory %s: %v", filepath.Dir(filePath), err)
		}
		if err := os.WriteFile(filePath, []byte(file.content), 0o600); err != nil {
			t.Fatalf("write fixture file %s: %v", filePath, err)
		}
	}
	return targetProjectPath(projectDirectory)
}

func (fixture targetProjectFixture) elementIDs(elementType oftElementType) []oftElementID {
	var ids []oftElementID
	for _, file := range fixture.files {
		for _, match := range oftElementIDPattern.FindAllStringSubmatch(file.content, -1) {
			if oftElementType(match[1]) == elementType {
				ids = append(ids, oftElementID(match[0]))
			}
		}
	}
	return ids
}

func sameElementIDs(expected, actual []oftElementID) bool {
	if len(expected) != len(actual) {
		return false
	}
	expectedCopy := append([]oftElementID(nil), expected...)
	actualCopy := append([]oftElementID(nil), actual...)
	sort.Slice(expectedCopy, func(left, right int) bool { return expectedCopy[left] < expectedCopy[right] })
	sort.Slice(actualCopy, func(left, right int) bool { return actualCopy[left] < actualCopy[right] })
	for index := range expectedCopy {
		if expectedCopy[index] != actualCopy[index] {
			return false
		}
	}
	return true
}

type scenarioSuite[Response any] struct {
	catalog  scenarioCatalogPath
	fixtures fixtureCatalogPath
	boundary publicBoundaryProgram
	decode   responseDecoder[Response]
	match    responseMatcher[Response]
	bindings map[scenarioID]scenarioBinding[Response]
}

func newScenarioSuite[Response any](catalog scenarioCatalogPath, fixtures fixtureCatalogPath, boundary publicBoundaryProgram, decode responseDecoder[Response], match responseMatcher[Response]) *scenarioSuite[Response] {
	return &scenarioSuite[Response]{
		catalog:  catalog,
		fixtures: fixtures,
		boundary: boundary,
		decode:   decode,
		match:    match,
		bindings: make(map[scenarioID]scenarioBinding[Response]),
	}
}

func (suite *scenarioSuite[Response]) bind(id scenarioID, expected Response) {
	if _, exists := suite.bindings[id]; exists {
		panic("duplicate acceptance binding for " + string(id))
	}
	suite.bindings[id] = scenarioBinding[Response]{expected: expected}
}

func (suite *scenarioSuite[Response]) run(t *testing.T) {
	t.Helper()
	ids, err := canonicalScenarioIDs(suite.catalog)
	if err != nil {
		t.Fatalf("read canonical scenarios: %v", err)
	}
	fixtures, err := readFixtures(suite.fixtures)
	if err != nil {
		t.Fatalf("read acceptance fixtures: %v", err)
	}
	for _, id := range ids {
		id := id
		t.Run(string(id), func(t *testing.T) {
			suite.runScenario(t, id, fixtures)
		})
	}
	for id := range fixtures {
		if !containsScenario(ids, id) {
			t.Errorf("%s: fixture has no canonical scenario", id)
		}
	}
}

func (suite *scenarioSuite[Response]) runScenario(t *testing.T, id scenarioID, fixtures fixtureCatalog) {
	t.Helper()
	if suite.boundary.source == "" {
		t.Fatalf("%s: missing public-boundary binding", id)
	}
	binding, exists := suite.bindings[id]
	if !exists {
		t.Fatalf("%s: missing public-boundary binding", id)
	}
	fixture, exists := fixtures[id]
	if !exists {
		t.Fatalf("%s: missing fixture", id)
	}
	path := fixture.create(t)
	output, err := suite.boundary.execute(context.Background(), path)
	if err != nil {
		t.Fatalf("%s: public-boundary binding failed: %v", id, err)
	}
	actual, err := suite.decode(output)
	if err != nil {
		t.Fatalf("%s: decode public-boundary response: %v", id, err)
	}
	if err := suite.match(binding.expected, actual, fixture); err != nil {
		t.Fatalf("%s: %v", id, err)
	}
}

func (program publicBoundaryProgram) execute(ctx context.Context, path targetProjectPath) (programOutput, error) {
	directory, err := os.MkdirTemp("", "coding-agent-flow-acceptance-binding-*")
	if err != nil {
		return "", fmt.Errorf("create binding directory: %w", err)
	}
	defer os.RemoveAll(directory)

	moduleFile := "module acceptancebinding\n\ngo 1.26.8\n\nrequire github.com/exasol/coding-agent-flow v0.0.0\n\nreplace github.com/exasol/coding-agent-flow => " + repositoryRoot() + "\n"
	if err := os.WriteFile(filepath.Join(directory, "go.mod"), []byte(moduleFile), 0o600); err != nil {
		return "", fmt.Errorf("write binding module: %w", err)
	}
	if err := os.WriteFile(filepath.Join(directory, "main.go"), []byte(program.source), 0o600); err != nil {
		return "", fmt.Errorf("write binding program: %w", err)
	}

	command := exec.CommandContext(ctx, "go", "run", "-mod=mod", ".", string(path))
	command.Dir = directory
	output, err := command.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message == "" {
			message = err.Error()
		}
		return "", errors.New(message)
	}
	return programOutput(strings.TrimSpace(string(output))), nil
}

func TestPublicBoundaryProgramReportsIncompatibleSignature(t *testing.T) {
	program := publicBoundaryProgram{source: `package main

import "context"

type validator func(context.Context, string) (bool, error)

func incompatibleOFTTracingValidator(context.Context, string) (string, error) {
	return "", nil
}

var _ validator = incompatibleOFTTracingValidator

func main() {}
`}
	_, err := program.execute(context.Background(), targetProjectPath(t.TempDir()))
	if err == nil {
		t.Fatal("expected incompatible public-boundary signature to fail")
	}
	t.Logf("incompatible public-boundary diagnostic: %v", err)
	if !strings.Contains(err.Error(), "cannot use incompatibleOFTTracingValidator") {
		t.Fatalf("expected incompatible-signature error, got: %v", err)
	}
}

type fixtureCatalog map[scenarioID]targetProjectFixture

func canonicalScenarioIDs(catalog scenarioCatalogPath) ([]scenarioID, error) {
	file, err := os.Open(string(catalog))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var scenarioIDs []scenarioID
	pendingScenarioID := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "@id:") {
			pendingScenarioID = strings.TrimPrefix(line, "@id:")
			continue
		}
		if strings.HasPrefix(line, "Scenario: ") {
			if pendingScenarioID == "" {
				return nil, fmt.Errorf("%s has no @id tag", line)
			}
			scenarioIDs = append(scenarioIDs, scenarioID(pendingScenarioID))
			pendingScenarioID = ""
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return scenarioIDs, nil
}

func containsScenario(scenarioIDs []scenarioID, id scenarioID) bool {
	for _, scenarioID := range scenarioIDs {
		if scenarioID == id {
			return true
		}
	}
	return false
}

func readFixtures(path fixtureCatalogPath) (fixtureCatalog, error) {
	file, err := os.Open(string(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fixtures := make(fixtureCatalog)
	var scenarioID scenarioID
	var filePath string
	var content []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.TrimSpace(line) == "" && filePath == "":
			continue
		case strings.HasPrefix(line, "scenario "):
			if scenarioID != "" || filePath != "" {
				return nil, fmt.Errorf("start scenario before ending previous scenario")
			}
			scenarioID = scenarioIDFrom(line)
		case strings.HasPrefix(line, "file "):
			if scenarioID == "" || filePath != "" {
				return nil, fmt.Errorf("start file outside a scenario")
			}
			filePath = strings.TrimPrefix(line, "file ")
		case line == "end file":
			if filePath == "" {
				return nil, fmt.Errorf("end file without a file")
			}
			fixture := fixtures[scenarioID]
			fixture.files = append(fixture.files, targetProjectFile{path: filePath, content: strings.Join(content, "\n")})
			fixtures[scenarioID] = fixture
			filePath = ""
			content = nil
		case line == "end scenario":
			if scenarioID == "" || filePath != "" {
				return nil, fmt.Errorf("end incomplete scenario")
			}
			scenarioID = ""
		default:
			if filePath == "" {
				return nil, fmt.Errorf("fixture content outside a file")
			}
			content = append(content, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if scenarioID != "" || filePath != "" {
		return nil, fmt.Errorf("fixture file ends before scenario completion")
	}
	return fixtures, nil
}

func scenarioIDFrom(line string) scenarioID {
	return scenarioID(strings.TrimPrefix(line, "scenario "))
}

func repositoryRoot() string {
	return filepath.Clean(filepath.Join(filepath.Dir(testSourceFile()), "..", ".."))
}

func testSourceFile() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("determine acceptance test location")
	}
	return file
}

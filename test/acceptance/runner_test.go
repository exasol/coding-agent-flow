package acceptance

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func writeAtData(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), atDataFile)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReadAtDataReportsAnUnknownName(t *testing.T) {
	path := writeAtData(t, "[[scenario]]\n\n[[scenario.given.file]]\npath = \"{{missing}}\"\n")

	if _, err := readAtData(path); err == nil {
		t.Fatal("expected an error for {{missing}}")
	}
}

func TestReadAtDataFillsGivenValuesIntoASharedDocument(t *testing.T) {
	path := writeAtData(t, `
[variables]
item_document = '''
Items
=====

{{item_id}}
'''

[[scenario]]

[scenario.given]
location = "items/item.txt"
item_id = "item-1"

[[scenario.given.file]]
path = "{{location}}"
content = "${item_document}"

[scenario.then]
found = true
item_id = "{{item_id}}"
`)

	runs, err := readAtData(path)

	if err != nil {
		t.Fatal(err)
	}
	run := runs[0]
	wantFiles := []inputFile{{
		Path:    "items/item.txt",
		Content: "Items\n=====\n\nitem-1\n",
	}}
	if !reflect.DeepEqual(run.files, wantFiles) {
		t.Errorf("got files %q, want %q", run.files, wantFiles)
	}
	if want := (Values{"found": "true", "item_id": "item-1"}); !reflect.DeepEqual(run.then, want) {
		t.Errorf("got then %v, want %v", run.then, want)
	}
}

func TestReadAtDataRefusesAnUnknownTable(t *testing.T) {
	path := writeAtData(t, "[[scenario]]\n\n[[scenario.file]]\npath = \"a.md\"\n")

	if _, err := readAtData(path); err == nil {
		t.Fatal("expected an error for [[scenario.file]], which belongs under [scenario.given]")
	}
}

func TestReadAtDataRefusesAFileWithoutBlocks(t *testing.T) {
	if _, err := readAtData(writeAtData(t, "")); err == nil {
		t.Fatal("expected an error for an at_data.toml without [[scenario]] blocks")
	}
}

func TestReadAtDataReportsAMissingFile(t *testing.T) {
	if _, err := readAtData(filepath.Join(t.TempDir(), atDataFile)); err == nil {
		t.Fatal("expected an error for a missing at_data.toml")
	}
}

func TestReadAtDataReportsInvalidTOML(t *testing.T) {
	if _, err := readAtData(writeAtData(t, "[[scenario]\n")); err == nil {
		t.Fatal("expected an error for invalid TOML")
	}
}

func TestReadAtDataReportsAnUnknownVariable(t *testing.T) {
	path := writeAtData(t, "[[scenario]]\n\n[[scenario.given.file]]\npath = \"a.md\"\ncontent = \"${missing}\"\n")

	if _, err := readAtData(path); err == nil || !strings.Contains(err.Error(), "${missing}") {
		t.Fatalf("expected an error naming ${missing}, got %v", err)
	}
}

func TestReadAtDataReportsAnUnknownNameInThen(t *testing.T) {
	path := writeAtData(t, "[[scenario]]\n\n[scenario.then]\nitem_id = \"{{item_id}}\"\n")

	if _, err := readAtData(path); err == nil || !strings.Contains(err.Error(), "{{item_id}}") {
		t.Fatalf("expected an error naming {{item_id}}, got %v", err)
	}
}

func TestReadAtDataReportsAFileWithoutPath(t *testing.T) {
	path := writeAtData(t, "[[scenario]]\n\n[[scenario.given.file]]\ncontent = \"text\"\n")

	if _, err := readAtData(path); err == nil || !strings.Contains(err.Error(), "has no path") {
		t.Fatalf("expected an error for a file without path, got %v", err)
	}
}

func TestReadAtDataRefusesASingleFileTable(t *testing.T) {
	path := writeAtData(t, "[[scenario]]\n\n[scenario.given.file]\npath = \"a.md\"\n")

	if _, err := readAtData(path); err == nil || !strings.Contains(err.Error(), "[[scenario.given.file]]") {
		t.Fatalf("expected an error for [scenario.given.file], got %v", err)
	}
}

func TestReadAtDataKeepsAThenValueNamedFile(t *testing.T) {
	runs, err := readAtData(writeAtData(t, "[[scenario]]\n\n[scenario.then]\nfile = \"item.txt\"\n"))

	if err != nil || runs[0].then.String("file") != "item.txt" {
		t.Fatalf("got %v, %v", runs, err)
	}
}

func TestReadAtDataNamesRunsByIDAndDescription(t *testing.T) {
	runs, err := readAtData(writeAtData(t, "[[scenario]]\nid = \"scn~list-items~1\"\ndescription = \"List the items\"\n\n[[scenario]]\n"))

	if err != nil {
		t.Fatal(err)
	}
	if runs[0].name() != "scn~list-items~1: List the items" || runs[1].name() != "run-2" {
		t.Fatalf("got %q and %q", runs[0].name(), runs[1].name())
	}
}

func TestReadAtDataRefusesAnEmptyFilePath(t *testing.T) {
	path := writeAtData(t, "[[scenario]]\n\n[[scenario.given.file]]\npath = \"\"\n")

	if _, err := readAtData(path); err == nil || !strings.Contains(err.Error(), "has no path") {
		t.Fatalf("expected an error for an empty path, got %v", err)
	}
}

func TestWriteFilesStaysInsideTheFileDirectory(t *testing.T) {
	run := testRun{files: []inputFile{{Path: "../escape.md"}}}
	if err := run.writeFiles(t.TempDir()); err == nil {
		t.Fatal("expected an error for a path outside the file directory")
	}
}

func TestWriteFilesWritesNestedFiles(t *testing.T) {
	directory := t.TempDir()
	run := testRun{files: []inputFile{{Path: "items/a/item.txt", Content: "item-1"}}}

	if err := run.writeFiles(directory); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(directory, "items", "a", "item.txt"))
	if err != nil || string(content) != "item-1" {
		t.Fatalf("got %q, %v", content, err)
	}
}

func TestWriteFilesReportsADirectoryThatIsAFile(t *testing.T) {
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, "items"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	run := testRun{files: []inputFile{{Path: "items/item.txt"}}}

	if err := run.writeFiles(directory); err == nil || !strings.Contains(err.Error(), "create input directory") {
		t.Fatalf("expected a directory error, got %v", err)
	}
}

func TestWriteFilesReportsAFileThatIsADirectory(t *testing.T) {
	directory := t.TempDir()
	if err := os.Mkdir(filepath.Join(directory, "item.txt"), 0o755); err != nil {
		t.Fatal(err)
	}
	run := testRun{files: []inputFile{{Path: "item.txt"}}}

	if err := run.writeFiles(directory); err == nil || !strings.Contains(err.Error(), "write input file") {
		t.Fatalf("expected a write error, got %v", err)
	}
}

func TestCallProgramReportsACompileError(t *testing.T) {
	_, err := newCallProgram(`strconv.Atoi(42)`, []string{`"strconv"`}, nil).build(t.Context(), t.TempDir())

	if err == nil || !strings.Contains(err.Error(), "cannot use") {
		t.Fatalf("expected a compile error, got %v", err)
	}
}

func TestCallProgramPassesTheFileDirectoryToTheCall(t *testing.T) {
	program, err := newCallProgram(`files.Directory(), error(nil)`, nil, nil).build(t.Context(), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	output, err := program.call(t.Context(), "/inputs", nil)

	if err != nil || string(output) != `"/inputs"` {
		t.Fatalf("got %s, %v", output, err)
	}
}

func TestCallProgramPassesGivenValuesToTheCall(t *testing.T) {
	program, err := newCallProgram(`given.String("item_id"), error(nil)`, nil, nil).build(t.Context(), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	output, err := program.call(t.Context(), "/inputs", Values{"item_id": "item-1"})

	if err != nil || string(output) != `"item-1"` {
		t.Fatalf("got %s, %v", output, err)
	}
}

func TestCallProgramReportsTheErrorOfTheCall(t *testing.T) {
	program, err := newCallProgram(`"", errors.New("call failed")`, []string{`"errors"`}, nil).build(t.Context(), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	if _, err := program.call(t.Context(), "", nil); err == nil || err.Error() != "call failed" {
		t.Fatalf("expected the call's error, got %v", err)
	}
}

func TestCallProgramReportsAnUnwritableDirectory(t *testing.T) {
	_, err := newCallProgram(`"", error(nil)`, nil, nil).build(t.Context(), filepath.Join(t.TempDir(), "missing"))

	if err == nil || !strings.Contains(err.Error(), "write") {
		t.Fatalf("expected a write error, got %v", err)
	}
}

func TestRunCommandReportsACommandThatCannotStart(t *testing.T) {
	if _, err := runCommand(t.Context(), "", filepath.Join(t.TempDir(), "missing-command")); err == nil {
		t.Fatal("expected an error for a missing command")
	}
}

func TestSuiteWithoutChecksIsRefused(t *testing.T) {
	if err := NewSuite(When[struct{}]{}).requireChecks(); err == nil {
		t.Fatal("expected an error for a suite without checks")
	}
}

func TestRunCheckReportsAMissingValue(t *testing.T) {
	check := Check[struct{}](func(_ struct{}, _ Values, then Values) error {
		return Expect("valid", then.Bool("valid"), true)
	})

	if err := runCheck(check, struct{}{}, Values{}, Values{}); err == nil || err.Error() != `at_data.toml has no value "valid"` {
		t.Fatalf("got %v", err)
	}
}

func TestRunCheckPassesGivenValues(t *testing.T) {
	check := Check[struct{}](func(_ struct{}, given Values, _ Values) error {
		return Expect("item ID", "item-1", given.String("item_id"))
	})

	if err := runCheck(check, struct{}{}, Values{"item_id": "item-1"}, Values{}); err != nil {
		t.Fatal(err)
	}
}

func TestExpectNamesTheDifference(t *testing.T) {
	if err := Expect("valid", true, false); err == nil || err.Error() != "valid: expected true, got false" {
		t.Fatalf("got %v", err)
	}
}

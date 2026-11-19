package openfasttrace

import (
	"strings"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

// [utest -> req~oft-requirements-detection~1]
func TestReadASpecReportReturnsTheOFTSpecificationItemIDsPerArtifactType(t *testing.T) {
	report := []byte(`<specdocument>
  <specobjects doctype="feat"><specobject><id>oft-validation</id><version>1</version></specobject></specobjects>
  <specobjects doctype="req">
    <specobject><id>first</id><version>1</version></specobject>
    <specobject><id>second</id><version>2</version></specobject>
  </specobjects>
</specdocument>`)

	validation, err := readASpecReport(report, true)

	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"feat~oft-validation~1", "req~first~1", "req~second~2"} {
		id, err := domain.ParseOFTSpecificationItemID(text)
		if err != nil {
			t.Fatal(err)
		}
		if !validation.Contains(id) {
			t.Errorf("expected %s in %+v", text, validation)
		}
	}
	if !validation.Valid() {
		t.Error("expected valid tracing")
	}
}

func TestReadASpecReportReportsAnUnreadableReport(t *testing.T) {
	if _, err := readASpecReport([]byte("Exception in thread main"), false); err == nil {
		t.Fatal("expected an error for output that is no ASpec report")
	}
}

func TestReadASpecReportRefusesAnotherRootElement(t *testing.T) {
	for _, output := range []string{
		`<html><body>OpenFastTrace</body></html>`,
		`<specobjects doctype="feat"><specobject><id>f</id><version>1</version></specobject></specobjects>`,
	} {
		if _, err := readASpecReport([]byte(output), true); err == nil || !strings.Contains(err.Error(), "specdocument") {
			t.Errorf("expected an error naming specdocument for %q, got %v", output, err)
		}
	}
}

func TestReadASpecReportReadsAnEmptyASpecDocument(t *testing.T) {
	validation, err := readASpecReport([]byte(`<?xml version="1.0" encoding="UTF-8"?><specdocument></specdocument>`), true)

	if err != nil || !validation.Valid() {
		t.Fatalf("expected a valid empty result, got %+v, %v", validation, err)
	}
}

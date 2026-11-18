package domain

import (
	"encoding/json"
	"reflect"
	"testing"
)

// [utest -> req~link-validation~1]
func TestNewLinkValidationSortsBrokenLinks(t *testing.T) {
	second := NewBrokenLink(NewLink("b.md", 1, "x.md"), "File not found")
	first := NewBrokenLink(NewLink("a.md", 9, "y.md"), "File not found")
	firstOnLine := NewBrokenLink(NewLink("a.md", 2, "z.md#nope"), "Cannot find fragment")

	validation := NewLinkValidation([]BrokenLink{second, first, firstOnLine})

	if want := []BrokenLink{firstOnLine, first, second}; !reflect.DeepEqual(validation.BrokenLinks(), want) {
		t.Fatalf("got %+v, want %+v", validation.BrokenLinks(), want)
	}
	if validation.Valid() {
		t.Fatal("expected invalid with Broken Links")
	}
	if !validation.ContainsBrokenLink("a.md", "z.md#nope") || validation.ContainsBrokenLink("b.md", "z.md#nope") {
		t.Fatal("ContainsBrokenLink must match source and target")
	}
}

func TestLinkValidationUsesThePublishedResponseFormat(t *testing.T) {
	validation := NewLinkValidation([]BrokenLink{NewBrokenLink(NewLink("README.md", 7, "doc/missing.md"), "File not found")})

	encoded, err := json.Marshal(validation)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"valid":false,"brokenLinks":[{"source":"README.md","line":7,"target":"doc/missing.md","status":"File not found"}]}`
	if string(encoded) != want {
		t.Fatalf("got %s, want %s", encoded, want)
	}
	var decoded LinkValidation
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, validation) {
		t.Fatalf("got %+v, want %+v", decoded, validation)
	}
}

func TestValidLinkValidationHasAnEmptyBrokenLinksList(t *testing.T) {
	encoded, err := json.Marshal(NewLinkValidation(nil))
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"valid":true,"brokenLinks":[]}`; string(encoded) != want {
		t.Fatalf("got %s, want %s", encoded, want)
	}
}

func TestZeroLinkValidationIsInvalid(t *testing.T) {
	var validation LinkValidation

	if validation.Valid() {
		t.Fatal("expected the zero LinkValidation to be invalid")
	}
	if !NewLinkValidation(nil).Valid() {
		t.Fatal("expected a LinkValidation without Broken Links to be valid")
	}
}

func TestLinkValidationRefusesInconsistentResponses(t *testing.T) {
	for _, data := range []string{
		`{"valid":true,"brokenLinks":[{"source":"a.md","line":1,"target":"b.md","status":"File not found"}]}`,
		`{"valid":false,"brokenLinks":[]}`,
	} {
		var validation LinkValidation
		if err := json.Unmarshal([]byte(data), &validation); err == nil {
			t.Errorf("expected %s to be refused", data)
		}
	}
}

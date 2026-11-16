package domain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func mustParse(t *testing.T, text string) OFTElementID {
	t.Helper()
	id, err := ParseOFTElementID(text)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func TestParseOFTElementIDReadsTypeNameAndRevision(t *testing.T) {
	id := mustParse(t, "req~oft-validation~2")

	if id.Type() != "req" || id.String() != "req~oft-validation~2" {
		t.Fatalf("got type %s, text %s", id.Type(), id)
	}
}

func TestParseOFTElementIDRefusesMalformedIDs(t *testing.T) {
	for _, text := range []string{"", "req~name", "req~name~one", "~name~1", "req~~1", "req~name~-1", "req~a~b~1"} {
		if _, err := ParseOFTElementID(text); err == nil {
			t.Errorf("expected %q to be refused", text)
		}
	}
}

func TestNewOFTValidationGroupsElementIDsByType(t *testing.T) {
	feature, requirement := mustParse(t, "feat~f~1"), mustParse(t, "req~r~1")

	validation := NewOFTValidation(true, []OFTElementID{feature, requirement})

	if !validation.Valid() || !validation.Contains(feature) || !validation.Contains(requirement) {
		t.Fatalf("got %+v", validation)
	}
	if validation.Contains(mustParse(t, "req~f~1")) {
		t.Fatal("expected req~f~1 to be absent: f is a feature name")
	}
}

func TestOFTValidationUsesThePublishedResponseFormat(t *testing.T) {
	validation := NewOFTValidation(false, []OFTElementID{mustParse(t, "req~r~1")})

	encoded, err := json.Marshal(validation)
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"valid":false,"oftElementIDs":{"req":["req~r~1"]}}`; string(encoded) != want {
		t.Fatalf("got %s, want %s", encoded, want)
	}
	var decoded OFTValidation
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, validation) {
		t.Fatalf("got %+v, want %+v", decoded, validation)
	}
}

func TestOFTValidationRefusesAnIDUnderAnotherType(t *testing.T) {
	var validation OFTValidation
	if err := json.Unmarshal([]byte(`{"valid":true,"oftElementIDs":{"feat":["req~r~1"]}}`), &validation); err == nil {
		t.Fatal("expected an ID listed under another type to be refused")
	}
}

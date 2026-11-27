package domain

import (
	"encoding/json"
	"testing"
)

func result(validation string, status ValidationStatus) ValidationResult {
	return NewValidationResult(validation, status, "CODE", "message")
}

// [utest -> req~validation-result~1]
func TestNewValidationReportReportsTheMostSevereStatus(t *testing.T) {
	for name, test := range map[string]struct {
		statuses []ValidationStatus
		want     ValidationStatus
	}{
		"no results":                     {nil, NotApplicable},
		"not applicable only":            {[]ValidationStatus{NotApplicable, NotApplicable}, NotApplicable},
		"pass over not applicable":       {[]ValidationStatus{NotApplicable, Pass}, Pass},
		"finding over pass":              {[]ValidationStatus{Pass, Finding, NotApplicable}, Finding},
		"tool failure over finding":      {[]ValidationStatus{Finding, ToolFailure, Pass}, ToolFailure},
		"unknown status as tool failure": {[]ValidationStatus{Pass, "unknown"}, ToolFailure},
	} {
		var results []ValidationResult
		for index, status := range test.statuses {
			results = append(results, result(string(rune('a'+index)), status))
		}

		if got := NewValidationReport(results).Status(); got != test.want {
			t.Errorf("%s: got %s, want %s", name, got, test.want)
		}
	}
}

func TestValidationReportKeepsTheSelectedOrder(t *testing.T) {
	report := NewValidationReport([]ValidationResult{result("second", Pass), result("first", Finding)})

	results := report.Results()
	if len(results) != 2 || results[0].Validation() != "second" || results[1].Validation() != "first" {
		t.Fatalf("got %+v", results)
	}
	if found, ok := report.Result("first"); !ok || found.Status() != Finding {
		t.Fatalf("got %+v, %t", found, ok)
	}
	if _, ok := report.Result("missing"); ok {
		t.Fatal("expected no result for an unselected validation")
	}
}

func TestZeroValidationReportIsAToolFailure(t *testing.T) {
	var report ValidationReport

	if report.Status() != ToolFailure {
		t.Fatalf("got %s", report.Status())
	}
}

func TestValidationReportUsesThePublishedResponseFormat(t *testing.T) {
	report := NewValidationReport([]ValidationResult{NewValidationResult("oft-tracing-validator", Pass, "CAF-OFT-TRACING-VALID", "valid")})

	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"status":"pass","results":[{"validation":"oft-tracing-validator","status":"pass","code":"CAF-OFT-TRACING-VALID","message":"valid"}]}`
	if string(encoded) != want {
		t.Fatalf("got %s, want %s", encoded, want)
	}
}

// A result with an unsupported status, or without a status, is a tool failure in the JSON response.
func TestUnsupportedResultStatusIsAToolFailure(t *testing.T) {
	report := NewValidationReport([]ValidationResult{NewValidationResult("check", ValidationStatus("passed"), "CODE", "message"), {}})

	output, err := json.Marshal(report)

	want := `{"status":"tool_failure","results":[{"validation":"check","status":"tool_failure","code":"CODE","message":"message"},{"validation":"","status":"tool_failure","code":"","message":""}]}`
	if err != nil || string(output) != want {
		t.Fatalf("got %s, %v; want %s", output, err, want)
	}
}

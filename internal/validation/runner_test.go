package validation

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// [utest~caf-validation-routing~1 -> req~caf-validation-routing~1]
func TestRunReportsPassForPresentRequirementsDocument(t *testing.T) {
	projectRoot := projectWithRequirements(t)

	report := Run(context.Background(), projectRoot, []string{RequirementsDocumentValidator})

	assert.Equal(t, StatusPass, report.Status)
	require.Len(t, report.Results, 1)
	assert.Equal(t, StatusPass, report.Results[0].Status)
}

func TestRunReportsFindingForMissingRequirementsDocument(t *testing.T) {
	report := Run(context.Background(), t.TempDir(), []string{RequirementsDocumentValidator})

	assert.Equal(t, StatusFinding, report.Status)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "CAF-REQUIREMENTS-DOCUMENT-MISSING", report.Results[0].Code)
}

func TestRunReportsToolFailureForUnknownValidator(t *testing.T) {
	report := Run(context.Background(), t.TempDir(), []string{"unknown"})

	assert.Equal(t, StatusToolFailure, report.Status)
	assert.Equal(t, "CAF-UNKNOWN-VALIDATOR", report.Results[0].Code)
}

func TestRunReportsToolFailureWhenCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	report := Run(ctx, t.TempDir(), []string{RequirementsDocumentValidator})

	assert.Equal(t, StatusToolFailure, report.Status)
	assert.Equal(t, "CAF-VALIDATION-CANCELLED", report.Results[0].Code)
}

func TestNewReportFailsClosedForUnknownStatus(t *testing.T) {
	report := NewReport([]Result{{Status: Status("unexpected")}})

	assert.Equal(t, StatusToolFailure, report.Status)
}

func TestDecideRequestsMoreWorkForFinding(t *testing.T) {
	report := NewReport([]Result{{Status: StatusFinding, Code: "CAF-TEST", Message: "Fix the finding."}})

	decision := Decide(report, false)

	assert.Equal(t, ContinueWork, decision.Action)
	assert.Equal(t, "CAF-TEST: Fix the finding.", decision.Reason)
}

func TestDecideAllowsStopAfterOneContinuation(t *testing.T) {
	report := NewReport([]Result{{Status: StatusFinding, Code: "CAF-TEST", Message: "Fix the finding."}})

	decision := Decide(report, true)

	assert.Equal(t, AllowStop, decision.Action)
}

func projectWithRequirements(t *testing.T) string {
	t.Helper()
	projectRoot := t.TempDir()
	requirementsDirectory := filepath.Join(projectRoot, "doc", "requirements")
	require.NoError(t, os.MkdirAll(requirementsDirectory, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(requirementsDirectory, "system.md"), nil, 0o600))
	return projectRoot
}

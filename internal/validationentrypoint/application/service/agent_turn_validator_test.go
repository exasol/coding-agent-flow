package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	exaerror "github.com/exasol/error-reporting-go"

	"github.com/exasol/coding-agent-flow/internal/global"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/domain"
)

func validateTurn(t *testing.T, ctx context.Context, reader fakeReader, continuationActive bool) (domain.TurnDecision, error) {
	t.Helper()
	runner, err := NewValidationRunner(reader,
		fakeValidator{name: "passing", status: domain.Pass},
		fakeValidator{name: "finding", status: domain.Finding},
	)
	if err != nil {
		t.Fatal(err)
	}
	return NewAgentTurnValidator(runner).Validate(ctx, domain.NewStopEvent(global.NewTargetProject("project"), continuationActive))
}

// [utest -> req~agent-turn-validation~1]
// [utest -> req~agent-continuation-limit~1]
func TestAgentTurnValidatorDecidesFromTheValidationReport(t *testing.T) {
	for name, test := range map[string]struct {
		selection          []string
		continuationActive bool
		want               domain.TurnAction
		reason             string
	}{
		"pass":                             {[]string{"passing"}, false, domain.EndTurn, ""},
		"finding":                          {[]string{"passing", "finding"}, false, domain.ContinueTurn, "finding: finding CODE"},
		"unknown validation":               {[]string{"unknown"}, false, domain.ContinueTurn, "unknown: tool_failure " + UnknownValidationCode},
		"finding with continuation active": {[]string{"finding"}, true, domain.EndTurn, "finding: finding CODE"},
	} {
		decision, err := validateTurn(t, context.Background(), fakeReader{selection: test.selection}, test.continuationActive)

		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if decision.Action() != test.want || !strings.Contains(decision.Reason(), test.reason) {
			t.Errorf("%s: got %s %q, want %s with %q", name, decision.Action(), decision.Reason(), test.want, test.reason)
		}
		if len(decision.Report().Results()) != len(test.selection) {
			t.Errorf("%s: got report %+v", name, decision.Report().Results())
		}
	}
}

// [utest -> req~agent-turn-validation~1]
// [utest -> req~agent-continuation-limit~1]
func TestAgentTurnValidatorNamesTheErrorCodeOfAValidationRunnerError(t *testing.T) {
	cause := exaerror.New("E-CAF-5").Message("Failed to read {{path}}.").Parameter("path", "/home/user/project/.caf.yaml")
	for continuationActive, want := range map[bool]domain.TurnAction{false: domain.ContinueTurn, true: domain.EndTurn} {
		decision, err := validateTurn(t, context.Background(), fakeReader{err: cause}, continuationActive)

		if err != nil {
			t.Fatal(err)
		}
		if decision.Action() != want || !strings.Contains(decision.Reason(), "E-CAF-5") {
			t.Errorf("continuation active %t: got %s %q, want %s with E-CAF-5", continuationActive, decision.Action(), decision.Reason(), want)
		}
		if strings.Contains(decision.Reason(), "/home/user") {
			t.Errorf("continuation active %t: the reason %q contains the cause", continuationActive, decision.Reason())
		}
	}
}

// [utest -> req~agent-turn-validation~1]
func TestAgentTurnValidatorContinuesOnAnErrorWithoutCode(t *testing.T) {
	decision, err := validateTurn(t, context.Background(), fakeReader{err: errors.New("broken")}, false)

	if err != nil || decision.Action() != domain.ContinueTurn || strings.Contains(decision.Reason(), "broken") {
		t.Fatalf("got %s %q, %v", decision.Action(), decision.Reason(), err)
	}
}

func TestAgentTurnValidatorReturnsTheErrorOfAnEndedContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := validateTurn(t, ctx, fakeReader{selection: []string{"passing"}}, false); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v", err)
	}
}

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

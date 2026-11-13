// Package cli implements the machine-facing CAF command contract.
package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/exasol/coding-agent-flow/internal/adapter"
	"github.com/exasol/coding-agent-flow/internal/config"
	"github.com/exasol/coding-agent-flow/internal/validation"
)

const (
	exitSuccess     = 0
	exitFinding     = 1
	exitToolFailure = 2
)

// Run executes CAF using injected streams and working-directory discovery.
func Run(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer, getwd func() (string, error)) int {
	// [impl~caf-machine-validation~1 -> req~caf-machine-validation~1]
	if len(args) == 0 || args[0] != "validate" {
		fmt.Fprintln(stderr, "Usage: caf validate [--agent codex]")
		return exitToolFailure
	}

	flags := flag.NewFlagSet("caf validate", flag.ContinueOnError)
	flags.SetOutput(stderr)
	agentName := flags.String("agent", "", "coding-agent protocol: codex")
	if err := flags.Parse(args[1:]); err != nil {
		return exitToolFailure
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "caf validate does not accept positional arguments")
		return exitToolFailure
	}

	if *agentName == "" {
		return runCI(ctx, stdout, stderr, getwd)
	}
	agent, err := adapter.ParseAgent(*agentName)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitToolFailure
	}
	return runAgent(ctx, agent, stdin, stdout, stderr)
}

func runAgent(ctx context.Context, agent adapter.Agent, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	invocation, err := adapter.DecodeStop(agent, stdin)
	if err != nil {
		return writeAgentToolFailure(agent, invocation.AlreadyContinued, stdout, stderr, "CAF-HOOK-INPUT-INVALID", "CAF could not read the Stop-hook input.")
	}

	project, err := config.Discover(invocation.WorkingDirectory)
	if err != nil {
		return writeAgentReport(agent, validation.NewReport([]validation.Result{{
			Validator: "configuration",
			Status:    validation.StatusToolFailure,
			Code:      "CAF-CONFIGURATION-UNAVAILABLE",
			Message:   "CAF could not discover a valid project configuration.",
		}}), invocation.AlreadyContinued, stdout, stderr)
	}

	report := validation.Run(ctx, project.Root, project.Config.Validators)
	return writeAgentReport(agent, report, invocation.AlreadyContinued, stdout, stderr)
}

func writeAgentToolFailure(agent adapter.Agent, alreadyContinued bool, stdout io.Writer, stderr io.Writer, code string, message string) int {
	report := validation.NewReport([]validation.Result{{
		Validator: "agent-adapter",
		Status:    validation.StatusToolFailure,
		Code:      code,
		Message:   message,
	}})
	return writeAgentReport(agent, report, alreadyContinued, stdout, stderr)
}

func writeAgentReport(agent adapter.Agent, report validation.Report, alreadyContinued bool, stdout io.Writer, stderr io.Writer) int {
	decision := validation.Decide(report, alreadyContinued)
	if err := adapter.EncodeStop(agent, decision, stdout); err != nil {
		fmt.Fprintln(stderr, "CAF could not encode the agent response")
		return exitToolFailure
	}
	// A valid provider decision must be consumed even when validation found a
	// problem, so hook mode exits zero after successfully writing protocol JSON.
	return exitSuccess
}

func runCI(ctx context.Context, stdout io.Writer, stderr io.Writer, getwd func() (string, error)) int {
	workingDirectory, err := getwd()
	if err != nil {
		return writeCIReport(toolFailureReport("CAF-WORKING-DIRECTORY-UNAVAILABLE", "CAF could not determine the working directory."), stdout, stderr)
	}
	project, err := config.Discover(workingDirectory)
	if err != nil {
		return writeCIReport(toolFailureReport("CAF-CONFIGURATION-UNAVAILABLE", "CAF could not discover a valid project configuration."), stdout, stderr)
	}
	return writeCIReport(validation.Run(ctx, project.Root, project.Config.Validators), stdout, stderr)
}

func toolFailureReport(code string, message string) validation.Report {
	return validation.NewReport([]validation.Result{{
		Validator: "caf",
		Status:    validation.StatusToolFailure,
		Code:      code,
		Message:   message,
	}})
}

func writeCIReport(report validation.Report, stdout io.Writer, stderr io.Writer) int {
	if err := json.NewEncoder(stdout).Encode(report); err != nil {
		fmt.Fprintln(stderr, "CAF could not encode the CI result")
		return exitToolFailure
	}
	switch report.Status {
	case validation.StatusPass, validation.StatusNotApplicable:
		return exitSuccess
	case validation.StatusFinding:
		return exitFinding
	default:
		return exitToolFailure
	}
}

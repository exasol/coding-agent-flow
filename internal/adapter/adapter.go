// Package adapter maps provider hook protocols to CAF's shared contracts.
package adapter

import (
	"fmt"
	"io"
	"strings"

	"github.com/exasol/coding-agent-flow/internal/event"
	"github.com/exasol/coding-agent-flow/internal/validation"
)

// Agent identifies an explicitly selected coding-agent protocol.
type Agent string

const (
	Codex Agent = "codex"
)

// ParseAgent validates the value of the --agent option.
func ParseAgent(value string) (Agent, error) {
	switch Agent(strings.ToLower(value)) {
	case Codex:
		return Codex, nil
	default:
		return "", fmt.Errorf("unsupported agent %q; expected codex", value)
	}
}

// DecodeStop converts the selected provider's Stop payload to a generic event.
func DecodeStop(agent Agent, input io.Reader) (event.Event, error) {
	// [impl~caf-agent-adapter~1 -> req~caf-agent-selection~1]
	switch agent {
	case Codex:
		return decodeCodexStop(input)
	default:
		return event.Event{}, fmt.Errorf("unsupported agent %q", agent)
	}
}

// EncodeStop writes only the response JSON accepted by the selected provider.
func EncodeStop(agent Agent, decision validation.Decision, output io.Writer) error {
	switch agent {
	case Codex:
		return encodeCodexStop(decision, output)
	default:
		return fmt.Errorf("unsupported agent %q", agent)
	}
}

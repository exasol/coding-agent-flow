package domain

import (
	"strings"
	"unicode/utf8"
)

// MaximumDiagnosticsSize limits the diagnostics of a ValidationResult in bytes.
const MaximumDiagnosticsSize = 8 << 10

// ValidationResult is the result of one Validation: its name, status, code, message, and internal diagnostics.
type ValidationResult struct {
	validation  string
	status      ValidationStatus
	code        string
	message     string
	diagnostics string
}

// NewValidationResult creates a ValidationResult. Any status other than pass, finding, tool_failure, or not_applicable
// becomes tool_failure.
func NewValidationResult(validation string, status ValidationStatus, code, message string) ValidationResult {
	return ValidationResult{validation: validation, status: status.supported(), code: code, message: message}
}

// Validation is the Validation Name, such as oft-tracing-validator.
func (result ValidationResult) Validation() string {
	return result.validation
}

// Status is the Validation Result Status; a result without a status is a tool failure.
func (result ValidationResult) Status() ValidationStatus {
	return result.status.supported()
}

// Code identifies the result for machines, such as CAF-UNKNOWN-VALIDATOR.
func (result ValidationResult) Code() string {
	return result.code
}

func (result ValidationResult) Message() string {
	return result.message
}

// WithDiagnostics returns the result with the first MaximumDiagnosticsSize bytes of diagnostics, such as a tool error.
// Diagnostics are internal: the JSON response holds the code and the message only.
func (result ValidationResult) WithDiagnostics(diagnostics string) ValidationResult {
	if len(diagnostics) > MaximumDiagnosticsSize {
		diagnostics = diagnostics[:MaximumDiagnosticsSize]
		for !utf8.ValidString(diagnostics) {
			diagnostics = diagnostics[:len(diagnostics)-1]
		}
	}
	result.diagnostics = strings.Clone(diagnostics)
	return result
}

func (result ValidationResult) Diagnostics() string {
	return result.diagnostics
}

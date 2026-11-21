package domain

// ValidationResult is the result of one Validation: its name, status, code, and message. It holds no raw tool output,
// so callers can print and log it.
type ValidationResult struct {
	validation string
	status     ValidationStatus
	code       string
	message    string
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

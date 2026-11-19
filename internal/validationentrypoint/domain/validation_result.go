package domain

// ValidationResult is the result of one Validation: its name, status, code, and message.
type ValidationResult struct {
	validation string
	status     ValidationStatus
	code       string
	message    string
}

func NewValidationResult(validation string, status ValidationStatus, code, message string) ValidationResult {
	return ValidationResult{validation: validation, status: status, code: code, message: message}
}

// Validation is the Validation Name, such as oft-tracing-validator.
func (result ValidationResult) Validation() string {
	return result.validation
}

func (result ValidationResult) Status() ValidationStatus {
	return result.status
}

// Code identifies the result for machines, such as CAF-UNKNOWN-VALIDATOR.
func (result ValidationResult) Code() string {
	return result.code
}

func (result ValidationResult) Message() string {
	return result.message
}

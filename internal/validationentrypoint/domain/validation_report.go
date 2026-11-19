package domain

import (
	"encoding/json"
	"slices"
)

// ValidationReport holds one Validation Result per selected Validation, in the selected order, and the Overall Result. The zero value reports tool_failure, so the report that accompanies an error is a tool failure.
type ValidationReport struct {
	status  ValidationStatus
	results []ValidationResult
}

// NewValidationReport creates a ValidationReport whose Overall Result is the most severe status of the results, or not_applicable for an empty list.
func NewValidationReport(results []ValidationResult) ValidationReport {
	// [impl -> req~validation-result~1]
	status := NotApplicable
	for _, result := range results {
		if result.status.severity() > status.severity() {
			status = result.status
		}
	}
	if status.severity() == ToolFailure.severity() {
		status = ToolFailure
	}
	return ValidationReport{status: status, results: slices.Clone(results)}
}

// Status is the Overall Result.
func (report ValidationReport) Status() ValidationStatus {
	if report.status == "" {
		return ToolFailure
	}
	return report.status
}

func (report ValidationReport) Results() []ValidationResult {
	return slices.Clone(report.results)
}

// Result returns the Validation Result of the Validation with the given name.
func (report ValidationReport) Result(validation string) (ValidationResult, bool) {
	index := slices.IndexFunc(report.results, func(result ValidationResult) bool { return result.validation == validation })
	if index < 0 {
		return ValidationResult{}, false
	}
	return report.results[index], true
}

// validationResultResponse is the JSON form of ValidationResult.
type validationResultResponse struct {
	Validation string           `json:"validation"`
	Status     ValidationStatus `json:"status"`
	Code       string           `json:"code"`
	Message    string           `json:"message"`
}

// validationReportResponse is the JSON form of ValidationReport.
type validationReportResponse struct {
	Status  ValidationStatus           `json:"status"`
	Results []validationResultResponse `json:"results"`
}

// MarshalJSON writes the fields status and results.
func (report ValidationReport) MarshalJSON() ([]byte, error) {
	results := make([]validationResultResponse, 0, len(report.results))
	for _, result := range report.results {
		results = append(results, validationResultResponse{Validation: result.validation, Status: result.status, Code: result.code, Message: result.message})
	}
	return json.Marshal(validationReportResponse{Status: report.Status(), Results: results})
}

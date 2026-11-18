package codingagentflow

import (
	"github.com/exasol/coding-agent-flow/internal/codevalidation/adapter/lychee"
	"github.com/exasol/coding-agent-flow/internal/codevalidation/application/service"
	"github.com/exasol/coding-agent-flow/internal/codevalidation/domain"
)

// LinkValidation is the result of the Link Checker. It lists the Broken Links of a Target Project.
type LinkValidation = domain.LinkValidation

// BrokenLink is a Link that failed the check, with the status that tells why.
type BrokenLink = domain.BrokenLink

// LinkChecker checks the Links in the .md, .rst, .txt, and .xml documents of a Target Project.
type LinkChecker = service.LinkChecker

// LinkCheckerOption configures a LinkChecker.
type LinkCheckerOption func(*linkCheckerSettings)

// linkCheckerSettings holds the settings that LinkCheckerOptions change.
type linkCheckerSettings struct {
	scope domain.LinkScope
}

// WithExternalLinks makes the LinkChecker also check External Links. This needs network access.
func WithExternalLinks() LinkCheckerOption {
	return func(settings *linkCheckerSettings) {
		settings.scope = domain.LocalAndExternalLinks
	}
}

// NewLinkChecker creates a LinkChecker that uses lychee. It checks Local Links offline by default. WithExternalLinks
// adds External Links.
func NewLinkChecker(options ...LinkCheckerOption) LinkChecker {
	settings := linkCheckerSettings{scope: domain.LocalLinks}
	for _, option := range options {
		option(&settings)
	}
	return service.NewLinkChecker(lychee.NewAdapter(), settings.scope)
}

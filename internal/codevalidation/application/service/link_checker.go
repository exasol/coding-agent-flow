// Package service holds one application service per use case.
package service

import (
	"context"

	"github.com/exasol/coding-agent-flow/internal/codevalidation/application/port"
	"github.com/exasol/coding-agent-flow/internal/codevalidation/domain"
	"github.com/exasol/coding-agent-flow/internal/global"
)

type LinkChecker struct {
	finder port.BrokenLinkFinder
	scope  domain.LinkScope
}

// NewLinkChecker creates a LinkChecker that checks the Links in scope.
func NewLinkChecker(finder port.BrokenLinkFinder, scope domain.LinkScope) LinkChecker {
	return LinkChecker{finder: finder, scope: scope}
}

func (checker LinkChecker) Check(ctx context.Context, targetProject global.TargetProject) (domain.LinkValidation, error) {
	// [impl -> req~link-validation~1]
	// TODO(#54): decide link validity in the domain from the lychee report.
	return checker.finder.FindBrokenLinks(ctx, targetProject, checker.scope)
}

// Package port holds the interfaces that the application services call.
package port

import (
	"context"

	"github.com/exasol/coding-agent-flow/internal/codevalidation/domain"
	"github.com/exasol/coding-agent-flow/internal/global"
)

// BrokenLinkFinder checks the Links of a Target Project in the given scope and returns the Broken Links.
type BrokenLinkFinder interface {
	FindBrokenLinks(ctx context.Context, targetProject global.TargetProject, scope domain.LinkScope) (domain.LinkValidation, error)
}

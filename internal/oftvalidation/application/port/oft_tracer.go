// Package port holds the interfaces that the application services call.
package port

import (
	"context"

	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

type OFTTracer interface {
	Trace(ctx context.Context, targetProject domain.TargetProject) (domain.OFTValidation, error)
}

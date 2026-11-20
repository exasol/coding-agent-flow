// Package port holds the interfaces that the application services call.
package port

import (
	"context"

	"github.com/exasol/coding-agent-flow/internal/global"
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

// OFTTracer traces a Target Project with OFT and returns the OFT Trace.
type OFTTracer interface {
	Trace(ctx context.Context, targetProject global.TargetProject) (domain.OFTTrace, error)
}

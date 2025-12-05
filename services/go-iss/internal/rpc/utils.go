package rpc

import (
	"context"

	"github.com/google/uuid"
)

func GetTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace_id").(string); ok && traceID != "" {
		return traceID
	}
	return uuid.New().String()
}

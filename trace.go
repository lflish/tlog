package tlog

import (
	"context"
	"github.com/google/uuid"
	"strings"
)

const TraceID = "traceId"

func newTraceID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func NewTraceIdCtx(id string) context.Context {
	if id == "" {
		return context.WithValue(context.Background(), TraceID, newTraceID())
	}
	return context.WithValue(context.Background(), TraceID, newTraceID())
}

func getTraceId(ctx context.Context) string {
	if ctx == nil {
		return "nil"
	}

	v, ok := ctx.Value(TraceID).(string)
	if !ok {
		return ""
	}
	return v
}

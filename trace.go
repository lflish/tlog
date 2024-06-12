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

func NewTraceIdCtx() context.Context {
	return context.WithValue(context.Background(), TraceID, newTraceID()[0:10])
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

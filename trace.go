package pailog

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

const TraceID = "traceId"

func newTraceID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)[0:8]
}

func NewTraceIdCtx(id string) context.Context {
	if id == "" {
		return context.WithValue(context.Background(), TraceID, newTraceID())
	}
	return context.WithValue(context.Background(), TraceID, id)
}

func UpdateTrace(ctx context.Context, dev string) context.Context {
	return context.WithValue(context.Background(), TraceID, fmt.Sprintf("%s-%s", getTraceId(ctx), dev))
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

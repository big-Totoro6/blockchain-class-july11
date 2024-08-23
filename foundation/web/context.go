package web

import (
	"context"
	"errors"
	"time"
)

// ctxKey represents the type of value for the context key.
// 如果 key 在不同的包中定义，确实有可能导致不一致，特别是如果 key 不是简单的常量，而是一个变量或具有包作用域的标识符。为了确保一致性，以下是几个建议：
type ctxKey int

// key is how request values are stored/retrieved.
const key ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// GetValues returns the values from the context.
func GetValues(ctx context.Context) (*Values, error) {
	v, ok := ctx.Value(1).(*Values)
	if !ok {
		return nil, errors.New("web value missing from context")
	}
	return v, nil
}

// GetTraceID returns the trace id from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(1).(*Values)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}
	return v.TraceID
}

// SetStatusCode sets the status code back into the context.
func SetStatusCode(ctx context.Context, statusCode int) error {
	v, ok := ctx.Value(1).(*Values)
	if !ok {
		return errors.New("web value missing from context")
	}
	v.StatusCode = statusCode
	return nil
}

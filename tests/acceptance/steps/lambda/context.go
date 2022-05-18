package lambda

import (
	"context"
)

// ContextKey defines a type to store the integration lambda session in context.Context.
type ContextKey string

var contextKey ContextKey = "lambdaSession"

func InitializeContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, &Session{})
}

// GetSession returns the integration lambda session stored in context.
// Note that the context should be previously initialized with InitializeContext function.
func GetSession(ctx context.Context) *Session {
	return ctx.Value(contextKey).(*Session)
}

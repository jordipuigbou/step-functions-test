package db

import (
	"context"
)

// ContextKey defines a type to store the integration DB session in context.Context.
type ContextKey string

var contextKey ContextKey = "dbSession"

func InitializeContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, &Session{})
}

// GetSession returns the integration DB session stored in context.
// Note that the context should be previously initialized with InitializeContext function.
func GetSession(ctx context.Context) *Session {
	return ctx.Value(contextKey).(*Session)
}

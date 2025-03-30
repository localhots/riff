package riff

import (
	"context"
)

type contextKey struct{}

func WithContext(ctx context.Context, fields ...Field) context.Context {
	return context.WithValue(ctx, contextKey{}, append(FromContext(ctx), fields...))
}

func FromContext(ctx context.Context) []Field {
	v, _ := ctx.Value(contextKey{}).([]Field)
	return v
}

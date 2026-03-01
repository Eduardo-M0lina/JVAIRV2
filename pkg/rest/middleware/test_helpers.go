package middleware

import "context"

// WithTestAbilities adds abilities to the context for testing purposes
func WithTestAbilities(ctx context.Context, abilities []string) context.Context {
	return context.WithValue(ctx, abilityKey{}, abilities)
}

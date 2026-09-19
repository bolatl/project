package context

import (
	"context"

	"github.com/bolatl/lenslocked/models"
)

type key string 

const (
	userKey key = "user"
)

// WithUser returns a new context with the given user value.
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User retrieves the user value from the context. 
// If the value is not present or is of a different type, it returns nil.
func User(ctx context.Context) *models.User {
	user, ok := ctx.Value(userKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
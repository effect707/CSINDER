package auth

import (
	"context"
	"net/http"

	"github.com/effect707/csinder/internal/model"
)

type contextKey string

const userContextKey contextKey = "user"

func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUser(r *http.Request) *model.User {
	u, _ := r.Context().Value(userContextKey).(*model.User)
	return u
}

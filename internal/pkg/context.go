package pkg

import (
	"context"

	"github.com/google/uuid"
)

type contextKey struct{}

var userIDKey contextKey

func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	user_id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return user_id, ok
}

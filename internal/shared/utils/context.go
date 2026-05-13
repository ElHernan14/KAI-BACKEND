package utils

import "context"

type ContextKey string

const (
	UserIDKey ContextKey = "user_id"
)

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

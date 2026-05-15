package appcontext

import "context"

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	authUserKey  contextKey = "auth_user"
)

type AuthUser struct {
	UserID string
	Email  string
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func GetRequestID(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDKey).(string)
	return requestID, ok
}

func WithAuthUser(ctx context.Context, user AuthUser) context.Context {
	return context.WithValue(ctx, authUserKey, user)
}

func GetAuthUser(ctx context.Context) (AuthUser, bool) {
	user, ok := ctx.Value(authUserKey).(AuthUser)
	return user, ok
}

func GetUserID(ctx context.Context) (string, bool) {
	user, ok := GetAuthUser(ctx)
	if !ok || user.UserID == "" {
		return "", false
	}

	return user.UserID, true
}

func GetUserEmail(ctx context.Context) (string, bool) {
	user, ok := GetAuthUser(ctx)
	if !ok || user.Email == "" {
		return "", false
	}

	return user.Email, true
}

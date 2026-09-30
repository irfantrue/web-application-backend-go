package ctxval

import "context"

const userIDKey contextKey = "user_id"

// GetUserID retrieves the user ID from the given context.
// It returns an empty string if no user ID is found or if the value is not a string.
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(userIDKey).(string); ok {
		return id
	}
	return ""
}

// SetUserID returns a new context with the specified user ID attached to it.
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

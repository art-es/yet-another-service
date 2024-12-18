package context

import "context"

type keyUserID struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, keyUserID{}, userID)
}

func UserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(keyUserID{}).(string)
	return userID, ok
}

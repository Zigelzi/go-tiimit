package auth

import "context"

type UserInfo struct {
	IsLoggedIn bool
}

type contextKey string

const userInfoKey contextKey = "user_session_id"

func GetUserInfo(ctx context.Context) UserInfo {
	if user, ok := ctx.Value(userInfoKey).(UserInfo); ok {
		return user
	}
	return UserInfo{IsLoggedIn: false}
}

func WithUserInfo(ctx context.Context, user UserInfo) context.Context {
	return context.WithValue(ctx, userInfoKey, user)
}

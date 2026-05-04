package http

import (
	"context"
	"time"
)

type authContextKey struct{}

// AuthInfo は認証済みリクエストのユーザーとセッション情報です。
type AuthInfo struct {
	UserID    string
	SessionID string
	JTI       string
	ExpiresAt time.Time
}

func withAuthInfo(ctx context.Context, info AuthInfo) context.Context {
	return context.WithValue(ctx, authContextKey{}, info)
}

// AuthInfoFromContext はリクエスト context から認証情報を取得します。
func AuthInfoFromContext(ctx context.Context) (AuthInfo, bool) {
	info, ok := ctx.Value(authContextKey{}).(AuthInfo)
	return info, ok
}

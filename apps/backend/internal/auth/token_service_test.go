package auth

import (
	"strings"
	"testing"
	"time"
)

func TestIssueSessionTokens(t *testing.T) {
	service, err := NewTokenService("local-test-secret")
	if err != nil {
		t.Fatalf("NewTokenService returned error: %v", err)
	}

	now := time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)
	tokens, err := service.IssueSessionTokens("test-user", now)
	if err != nil {
		t.Fatalf("IssueSessionTokens returned error: %v", err)
	}

	if !strings.HasPrefix(tokens.SessionID, "session_") {
		t.Fatalf("SessionID = %q, want session_ prefix", tokens.SessionID)
	}
	if !strings.HasPrefix(tokens.JTI, "jti_") {
		t.Fatalf("JTI = %q, want jti_ prefix", tokens.JTI)
	}
	if len(strings.Split(tokens.AccessToken, ".")) != 3 {
		t.Fatalf("AccessToken = %q, want JWT format", tokens.AccessToken)
	}
	if !strings.HasPrefix(tokens.RefreshToken, "refresh_") {
		t.Fatalf("RefreshToken = %q, want refresh_ prefix", tokens.RefreshToken)
	}
	if tokens.AccessTokenEncrypted == tokens.AccessToken {
		t.Fatal("AccessTokenEncrypted must not equal AccessToken")
	}
	if tokens.RefreshTokenEncrypted == tokens.RefreshToken {
		t.Fatal("RefreshTokenEncrypted must not equal RefreshToken")
	}
	if !tokens.AccessTokenExpiresAt.Equal(now.Add(accessTokenTTL)) {
		t.Fatalf("AccessTokenExpiresAt = %v, want %v", tokens.AccessTokenExpiresAt, now.Add(accessTokenTTL))
	}
	if !tokens.RefreshTokenExpiresAt.Equal(now.Add(refreshTokenTTL)) {
		t.Fatalf("RefreshTokenExpiresAt = %v, want %v", tokens.RefreshTokenExpiresAt, now.Add(refreshTokenTTL))
	}
	if tokens.ExpiresIn != int32(accessTokenTTL/time.Second) {
		t.Fatalf("ExpiresIn = %d, want %d", tokens.ExpiresIn, int32(accessTokenTTL/time.Second))
	}
}

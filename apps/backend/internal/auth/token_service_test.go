package auth

import (
	"errors"
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
	if len(strings.Split(tokens.InternalAccessToken, ".")) != 3 {
		t.Fatalf("InternalAccessToken = %q, want JWT format", tokens.InternalAccessToken)
	}
	if !strings.HasPrefix(tokens.ProviderAccessToken, "provider_access_") {
		t.Fatalf("ProviderAccessToken = %q, want provider_access_ prefix", tokens.ProviderAccessToken)
	}
	if !strings.HasPrefix(tokens.ProviderRefreshToken, "provider_refresh_") {
		t.Fatalf("ProviderRefreshToken = %q, want provider_refresh_ prefix", tokens.ProviderRefreshToken)
	}
	if tokens.ProviderAccessTokenEncrypted == tokens.ProviderAccessToken {
		t.Fatal("ProviderAccessTokenEncrypted must not equal ProviderAccessToken")
	}
	if tokens.ProviderRefreshTokenEncrypted == tokens.ProviderRefreshToken {
		t.Fatal("ProviderRefreshTokenEncrypted must not equal ProviderRefreshToken")
	}
	if !tokens.InternalAccessTokenExpiresAt.Equal(now.Add(accessTokenTTL)) {
		t.Fatalf("InternalAccessTokenExpiresAt = %v, want %v", tokens.InternalAccessTokenExpiresAt, now.Add(accessTokenTTL))
	}
	if !tokens.ProviderAccessTokenExpiresAt.Equal(now.Add(accessTokenTTL)) {
		t.Fatalf("ProviderAccessTokenExpiresAt = %v, want %v", tokens.ProviderAccessTokenExpiresAt, now.Add(accessTokenTTL))
	}
	if !tokens.ProviderRefreshTokenExpiresAt.Equal(now.Add(refreshTokenTTL)) {
		t.Fatalf("ProviderRefreshTokenExpiresAt = %v, want %v", tokens.ProviderRefreshTokenExpiresAt, now.Add(refreshTokenTTL))
	}
	if tokens.InternalAccessTokenExpiresInSec != int32(accessTokenTTL/time.Second) {
		t.Fatalf("InternalAccessTokenExpiresInSec = %d, want %d", tokens.InternalAccessTokenExpiresInSec, int32(accessTokenTTL/time.Second))
	}

	claims, err := service.VerifyInternalJWT(tokens.InternalAccessToken, now)
	if err != nil {
		t.Fatalf("VerifyInternalJWT returned error: %v", err)
	}
	if claims.Subject != "test-user" {
		t.Fatalf("claims.Subject = %q, want test-user", claims.Subject)
	}
	if claims.SessionID != tokens.SessionID {
		t.Fatalf("claims.SessionID = %q, want %q", claims.SessionID, tokens.SessionID)
	}
	if claims.JTI != tokens.JTI {
		t.Fatalf("claims.JTI = %q, want %q", claims.JTI, tokens.JTI)
	}
}

func TestVerifyInternalJWTRejectsTamperedToken(t *testing.T) {
	service, err := NewTokenService("local-test-secret")
	if err != nil {
		t.Fatalf("NewTokenService returned error: %v", err)
	}

	now := time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)
	tokens, err := service.IssueSessionTokens("test-user", now)
	if err != nil {
		t.Fatalf("IssueSessionTokens returned error: %v", err)
	}

	tamperedToken := tokens.InternalAccessToken[:len(tokens.InternalAccessToken)-1] + "x"
	if _, err := service.VerifyInternalJWT(tamperedToken, now); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("VerifyInternalJWT error = %v, want ErrInvalidToken", err)
	}
}

func TestVerifyInternalJWTRejectsExpiredToken(t *testing.T) {
	service, err := NewTokenService("local-test-secret")
	if err != nil {
		t.Fatalf("NewTokenService returned error: %v", err)
	}

	now := time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)
	tokens, err := service.IssueSessionTokens("test-user", now)
	if err != nil {
		t.Fatalf("IssueSessionTokens returned error: %v", err)
	}

	expiredAt := now.Add(accessTokenTTL)
	if _, err := service.VerifyInternalJWT(tokens.InternalAccessToken, expiredAt); !errors.Is(err, ErrExpiredToken) {
		t.Fatalf("VerifyInternalJWT error = %v, want ErrExpiredToken", err)
	}
}

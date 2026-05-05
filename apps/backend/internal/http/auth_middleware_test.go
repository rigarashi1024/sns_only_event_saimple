package http

import (
	"net/http"
	"testing"
	"time"

	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/auth"
	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/repository"
)

func TestIsPublicEndpoint(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
		want bool
	}{
		{
			name: "healthz is public",
			req:  mustNewRequest(t, http.MethodGet, "/api/v1/healthz"),
			want: true,
		},
		{
			name: "login is public",
			req:  mustNewRequest(t, http.MethodPost, "/api/v1/auth/login"),
			want: true,
		},
		{
			name: "other paths are private",
			req:  mustNewRequest(t, http.MethodGet, "/auth/me"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPublicEndpoint(tt.req); got != tt.want {
				t.Fatalf("isPublicEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSessionValidForClaims(t *testing.T) {
	now := time.Date(2026, 5, 4, 12, 0, 0, 123, time.UTC)
	expiresAt := now.Add(15 * time.Minute)
	claims := &auth.InternalJWTClaims{
		Subject:   "test-user",
		SessionID: "session-001",
		JTI:       "jti-001",
		ExpiresAt: expiresAt.Unix(),
	}
	session := &repository.Session{
		UserID:                       "test-user",
		InternalJWTJTI:               "jti-001",
		InternalAccessTokenExpiresAt: expiresAt,
	}

	if !isSessionValidForClaims(session, claims, now) {
		t.Fatal("isSessionValidForClaims() = false, want true")
	}

	session.InternalJWTJTI = "different-jti"
	if isSessionValidForClaims(session, claims, now) {
		t.Fatal("isSessionValidForClaims() = true for mismatched jti, want false")
	}
}

func mustNewRequest(t *testing.T, method string, path string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatalf("http.NewRequest returned error: %v", err)
	}
	return req
}

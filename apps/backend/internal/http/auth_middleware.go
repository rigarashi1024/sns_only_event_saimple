package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/gen"
	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/auth"
	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/repository"
)

const internalAccessTokenCookieName = "internal_access_token"

// WithAuth は公開エンドポイント以外に HttpOnly cookie ベースの認証を適用します。
func (h *Handler) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicEndpoint(r) {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(internalAccessTokenCookieName)
		if err != nil || cookie.Value == "" {
			writeJSON(w, http.StatusUnauthorized, gen.ErrorResponse{Message: "authentication required"})
			return
		}

		now := time.Now().UTC()
		claims, err := h.tokenService.VerifyInternalJWT(cookie.Value, now)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, gen.ErrorResponse{Message: "invalid or expired token"})
			return
		}

		session, err := h.sessionRepo.FindByID(r.Context(), claims.SessionID)
		if err != nil {
			if errors.Is(err, repository.ErrSessionNotFound) {
				writeJSON(w, http.StatusUnauthorized, gen.ErrorResponse{Message: "invalid session"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, gen.ErrorResponse{Message: "failed to load session"})
			return
		}

		if !isSessionValidForClaims(session, claims, now) {
			writeJSON(w, http.StatusUnauthorized, gen.ErrorResponse{Message: "invalid session"})
			return
		}

		authInfo := AuthInfo{
			UserID:    claims.Subject,
			SessionID: claims.SessionID,
			JTI:       claims.JTI,
			ExpiresAt: time.Unix(claims.ExpiresAt, 0).UTC(),
		}
		next.ServeHTTP(w, r.WithContext(withAuthInfo(r.Context(), authInfo)))
	})
}

func isPublicEndpoint(r *http.Request) bool {
	return (r.Method == http.MethodGet && r.URL.Path == "/healthz") ||
		(r.Method == http.MethodPost && r.URL.Path == "/auth/login")
}

func isSessionValidForClaims(session *repository.Session, claims *auth.InternalJWTClaims, now time.Time) bool {
	if session.UserID != claims.Subject {
		return false
	}
	if session.InternalJWTJTI != claims.JTI {
		return false
	}
	if session.InternalAccessTokenExpiresAt.Unix() != claims.ExpiresAt {
		return false
	}
	if !now.Before(session.InternalAccessTokenExpiresAt) {
		return false
	}
	return true
}

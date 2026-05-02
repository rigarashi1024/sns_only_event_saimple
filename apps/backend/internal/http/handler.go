package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	gofirestore "cloud.google.com/go/firestore"
	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/gen"
	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/repository"
)

type Handler struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewHandler(firestoreClient *gofirestore.Client) *Handler {
	return &Handler{
		userRepo:    repository.NewUserRepository(firestoreClient),
		sessionRepo: repository.NewSessionRepository(firestoreClient),
	}
}

func (h *Handler) GetHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, gen.HealthzResponse{
		Status: "ok",
	})
}

func (h *Handler) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	var req gen.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, gen.ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	if req.LoginId == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, gen.ErrorResponse{
			Message: "login_id and password are required",
		})
		return
	}

	ctx := r.Context()

	user, err := h.userRepo.FindByLoginID(ctx, req.LoginId)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			writeJSON(w, http.StatusUnauthorized, gen.ErrorResponse{
				Message: "invalid credentials",
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, gen.ErrorResponse{
			Message: "failed to load user",
		})
		return
	}

	if req.Password != "password" {
		writeJSON(w, http.StatusUnauthorized, gen.ErrorResponse{
			Message: "invalid credentials",
		})
		return
	}

	now := time.Now().UTC()
	if err := h.sessionRepo.Create(ctx, repository.Session{
		UserID:                user.ID,
		AccessTokenEncrypted:  "dummy-access-token",
		AccessTokenExpiresAt:  now.Add(15 * time.Minute),
		RefreshTokenEncrypted: "dummy-refresh-token",
		RefreshTokenExpiresAt: now.Add(24 * time.Hour),
		InternalJWTJTI:        "dummy-jti",
		ProviderType:          "local",
	}); err != nil {
		writeJSON(w, http.StatusInternalServerError, gen.ErrorResponse{
			Message: "failed to create session",
		})
		return
	}

	writeJSON(w, http.StatusOK, gen.LoginResponse{
		AccessToken:  "dummy-access-token",
		RefreshToken: "dummy-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    900,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

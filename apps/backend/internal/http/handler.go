package http

import (
	"encoding/json"
	"net/http"

	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/gen"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
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

	if req.LoginId != "test-user" || req.Password != "password" {
		writeJSON(w, http.StatusUnauthorized, gen.ErrorResponse{
			Message: "invalid credentials",
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

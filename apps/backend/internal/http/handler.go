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

// Handler は OpenAPI 生成コードから呼ばれる HTTP ハンドラの実装です。
type Handler struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

// NewHandler は Firestore client を利用する repository を組み立てます。
func NewHandler(firestoreClient *gofirestore.Client) *Handler {
	return &Handler{
		userRepo:    repository.NewUserRepository(firestoreClient),
		sessionRepo: repository.NewSessionRepository(firestoreClient),
	}
}

// GetHealthz はアプリケーションの疎通確認用エンドポイントです。
func (h *Handler) GetHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, gen.HealthzResponse{
		Status: "ok",
	})
}

// PostAuthLogin はダミーログイン API です。
// 現時点では Firestore 上のユーザー存在確認と固定パスワードで判定します。
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

	// login_id は users/{login_id} の document id または email として扱う。
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

	// ダミー実装のため、パスワードは固定文字列のみ許可する。
	// 本実装では Firestore のハッシュ値検証などに置き換える想定。
	if req.Password != "password" {
		writeJSON(w, http.StatusUnauthorized, gen.ErrorResponse{
			Message: "invalid credentials",
		})
		return
	}

	now := time.Now().UTC()
	// ログイン成功時は sessions コレクションにセッション情報を残す。
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

// writeJSON は API レスポンスを JSON として返す共通処理です。
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

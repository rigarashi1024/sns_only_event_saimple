package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	gofirestore "cloud.google.com/go/firestore"
	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/gen"
	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/auth"
	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/repository"
)

// Handler は OpenAPI 生成コードから呼ばれる HTTP ハンドラの実装です。
type Handler struct {
	userRepo     *repository.UserRepository
	postRepo     *repository.PostRepository
	sessionRepo  *repository.SessionRepository
	timelineRepo *repository.TimelineRepository
	tokenService *auth.TokenService
	cookieSecure bool
}

// NewHandler は Firestore client を利用する repository を組み立てます。
func NewHandler(firestoreClient *gofirestore.Client, tokenService *auth.TokenService, cookieSecure bool) *Handler {
	return &Handler{
		userRepo:     repository.NewUserRepository(firestoreClient),
		postRepo:     repository.NewPostRepository(firestoreClient),
		sessionRepo:  repository.NewSessionRepository(firestoreClient),
		timelineRepo: repository.NewTimelineRepository(firestoreClient),
		tokenService: tokenService,
		cookieSecure: cookieSecure,
	}
}

// GetHealthz はアプリケーションの疎通確認用エンドポイントです。
func (h *Handler) GetHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, gen.HealthzResponse{
		Status: "ok",
	})
}

// GetProfile は認証済みユーザーのプロフィール情報と過去投稿を返します。
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	authInfo, ok := AuthInfoFromContext(r.Context())
	if !ok || authInfo.UserID == "" {
		writeJSON(w, http.StatusUnauthorized, gen.ErrorResponse{
			Message: "invalid authentication context",
		})
		return
	}

	user, err := h.userRepo.FindByID(r.Context(), authInfo.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			writeJSON(w, http.StatusInternalServerError, gen.ErrorResponse{
				Message: "authenticated user not found",
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, gen.ErrorResponse{
			Message: "failed to load user profile",
		})
		return
	}

	posts, err := h.postRepo.ListByUserID(r.Context(), authInfo.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, gen.ErrorResponse{
			Message: "failed to load user posts",
		})
		return
	}

	postItems := make([]gen.UserPostSummary, 0, len(posts))
	for _, post := range posts {
		postItems = append(postItems, gen.UserPostSummary{
			Id:           post.ID,
			UserId:       post.UserID,
			Content:      post.Content,
			LikeCount:    int32(post.LikeCount),
			CommentCount: int32(post.CommentCount),
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
		})
	}

	writeJSON(w, http.StatusOK, gen.ProfileResponse{
		User: gen.UserProfile{
			Id:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Nickname: user.Nickname,
		},
		Posts: postItems,
	})
}

// GetTimeline は認証済みユーザー向けに事前構築済み timelines を返します。
func (h *Handler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	authInfo, ok := AuthInfoFromContext(r.Context())
	if !ok || authInfo.UserID == "" {
		writeJSON(w, http.StatusUnauthorized, gen.ErrorResponse{
			Message: "invalid authentication context",
		})
		return
	}

	timelines, err := h.timelineRepo.ListByOwnerUserID(r.Context(), authInfo.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, gen.ErrorResponse{
			Message: "failed to load timeline",
		})
		return
	}

	items := make([]gen.TimelineEntry, 0, len(timelines))
	for _, timeline := range timelines {
		items = append(items, gen.TimelineEntry{
			Id:                  timeline.ID,
			PostId:              timeline.PostID,
			PostUserId:          timeline.PostUserID,
			Content:             timeline.Content,
			PostCreatedAt:       timeline.PostCreatedAt,
			TimelineOwnerUserId: timeline.TimelineOwnerUserID,
			CreatedAt:           timeline.CreatedAt,
			UpdatedAt:           timeline.UpdatedAt,
		})
	}

	writeJSON(w, http.StatusOK, gen.TimelineListResponse{
		Items: items,
	})
}

// PostAuthLogin はダミーログイン API です。
// Firestore 上のユーザー存在確認と bcrypt ハッシュ化済みパスワードで判定します。
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

	// Firestore に保存された bcrypt ハッシュとリクエストの平文パスワードを照合する。
	if user.PasswordHash == "" || !auth.VerifyPassword(req.Password, user.PasswordHash) {
		writeJSON(w, http.StatusUnauthorized, gen.ErrorResponse{
			Message: "invalid credentials",
		})
		return
	}

	now := time.Now().UTC()
	// 認証成功後、アプリ内部用の session_id / jti / access JWT / refresh token を発行する。
	tokens, err := h.tokenService.IssueSessionTokens(user.ID, now)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, gen.ErrorResponse{
			Message: "failed to issue tokens",
		})
		return
	}

	// sessions には provider token の暗号化済み値と、internal JWT 検証に必要なメタデータを保存する。
	if err := h.sessionRepo.Create(ctx, repository.Session{
		ID:                            tokens.SessionID,
		UserID:                        user.ID,
		ProviderType:                  "local",
		InternalJWTJTI:                tokens.JTI,
		InternalAccessTokenExpiresAt:  tokens.InternalAccessTokenExpiresAt,
		ProviderAccessTokenEncrypted:  tokens.ProviderAccessTokenEncrypted,
		ProviderAccessTokenExpiresAt:  tokens.ProviderAccessTokenExpiresAt,
		ProviderRefreshTokenEncrypted: tokens.ProviderRefreshTokenEncrypted,
		ProviderRefreshTokenExpiresAt: tokens.ProviderRefreshTokenExpiresAt,
	}); err != nil {
		writeJSON(w, http.StatusInternalServerError, gen.ErrorResponse{
			Message: "failed to create session",
		})
		return
	}

	// internal JWT はレスポンス body には載せず、JavaScript から読めない HttpOnly cookie として返す。
	http.SetCookie(w, &http.Cookie{
		Name:     internalAccessTokenCookieName,
		Value:    tokens.InternalAccessToken,
		Path:     "/",
		Expires:  tokens.InternalAccessTokenExpiresAt,
		MaxAge:   int(tokens.InternalAccessTokenExpiresInSec),
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteStrictMode,
	})

	writeJSON(w, http.StatusOK, gen.LoginResponse{
		InternalTokenType: "Bearer",
		InternalExpiresIn: tokens.InternalAccessTokenExpiresInSec,
	})
}

// writeJSON は API レスポンスを JSON として返す共通処理です。
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

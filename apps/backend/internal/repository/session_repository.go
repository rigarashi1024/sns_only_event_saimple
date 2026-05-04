package repository

import (
	"context"
	"errors"
	"time"

	gofirestore "cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrSessionNotFound は指定されたログインセッションが存在しない場合に返します。
var ErrSessionNotFound = errors.New("session not found")

// Session は Firestore の sessions コレクションに保存するログインセッション情報です。
type Session struct {
	ID                            string    `firestore:"id"`
	UserID                        string    `firestore:"user_id"`
	ProviderType                  string    `firestore:"provider_type"`
	InternalJWTJTI                string    `firestore:"internal_jwt_jti"`
	InternalAccessTokenExpiresAt  time.Time `firestore:"internal_access_token_expires_at"`
	ProviderAccessTokenEncrypted  string    `firestore:"provider_access_token_encrypted"`
	ProviderAccessTokenExpiresAt  time.Time `firestore:"provider_access_token_expires_at"`
	ProviderRefreshTokenEncrypted string    `firestore:"provider_refresh_token_encrypted"`
	ProviderRefreshTokenExpiresAt time.Time `firestore:"provider_refresh_token_expires_at"`
	CreatedAt                     time.Time `firestore:"created_at"`
	UpdatedAt                     time.Time `firestore:"updated_at"`
}

// SessionRepository は sessions コレクションへのアクセスをまとめます。
type SessionRepository struct {
	client *gofirestore.Client
}

// NewSessionRepository は Firestore client を使って SessionRepository を生成します。
func NewSessionRepository(client *gofirestore.Client) *SessionRepository {
	return &SessionRepository{client: client}
}

// Create はセッション ID と作成日時を補完して sessions コレクションへ保存します。
// 同じ ID の既存セッションがある場合は上書きせず、Firestore の AlreadyExists エラーを返します。
func (r *SessionRepository) Create(ctx context.Context, session Session) error {
	docRef := r.client.Collection("sessions").Doc(session.ID)
	// 呼び出し側が ID を指定しない場合は、Firestore の自動 ID で衝突を避ける。
	if session.ID == "" {
		docRef = r.client.Collection("sessions").NewDoc()
		session.ID = docRef.ID
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now().UTC()
	}
	if session.UpdatedAt.IsZero() {
		session.UpdatedAt = session.CreatedAt
	}

	_, err := docRef.Create(ctx, session)
	return err
}

// FindByID は sessions/{sessionID} を取得します。
func (r *SessionRepository) FindByID(ctx context.Context, sessionID string) (*Session, error) {
	doc, err := r.client.Collection("sessions").Doc(sessionID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	var session Session
	if err := doc.DataTo(&session); err != nil {
		return nil, err
	}
	session.ID = doc.Ref.ID
	return &session, nil
}

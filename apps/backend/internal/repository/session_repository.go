package repository

import (
	"context"
	"fmt"
	"time"

	gofirestore "cloud.google.com/go/firestore"
)

// Session は Firestore の sessions コレクションに保存するログインセッション情報です。
type Session struct {
	ID                    string    `firestore:"id"`
	UserID                string    `firestore:"user_id"`
	AccessTokenEncrypted  string    `firestore:"access_token_encrypted"`
	AccessTokenExpiresAt  time.Time `firestore:"access_token_expires_at"`
	RefreshTokenEncrypted string    `firestore:"refresh_token_encrypted"`
	RefreshTokenExpiresAt time.Time `firestore:"refresh_token_expires_at"`
	InternalJWTJTI        string    `firestore:"internal_jwt_jti"`
	ProviderType          string    `firestore:"provider_type"`
	CreatedAt             time.Time `firestore:"created_at"`
	UpdatedAt             time.Time `firestore:"updated_at"`
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
func (r *SessionRepository) Create(ctx context.Context, session Session) error {
	// 呼び出し側が ID を指定しない場合は、ユーザー ID と現在時刻から一意な ID を作る。
	if session.ID == "" {
		session.ID = fmt.Sprintf("session-%s-%d", session.UserID, time.Now().UTC().UnixNano())
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now().UTC()
	}
	if session.UpdatedAt.IsZero() {
		session.UpdatedAt = session.CreatedAt
	}

	_, err := r.client.Collection("sessions").Doc(session.ID).Set(ctx, session)
	return err
}

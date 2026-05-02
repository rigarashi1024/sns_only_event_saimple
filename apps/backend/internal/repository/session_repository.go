package repository

import (
	"context"
	"fmt"
	"time"

	gofirestore "cloud.google.com/go/firestore"
)

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

type SessionRepository struct {
	client *gofirestore.Client
}

func NewSessionRepository(client *gofirestore.Client) *SessionRepository {
	return &SessionRepository{client: client}
}

func (r *SessionRepository) Create(ctx context.Context, session Session) error {
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

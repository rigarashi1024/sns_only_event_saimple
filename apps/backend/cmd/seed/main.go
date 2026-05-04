package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/auth"
)

const defaultProjectID = "sns-only-event-local"
const seedPassword = "password"
const testUserID = "test-user"

type UserSeed struct {
	ID           string    `firestore:"id"`
	Name         string    `firestore:"name"`
	Email        string    `firestore:"email"`
	Nickname     string    `firestore:"nickname"`
	PasswordHash string    `firestore:"password_hash"`
	CreatedAt    time.Time `firestore:"created_at"`
	UpdatedAt    time.Time `firestore:"updated_at"`
}

type SessionSeed struct {
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

func main() {
	ctx := context.Background()

	projectID := os.Getenv("PUBSUB_PROJECT_ID")
	if projectID == "" {
		projectID = defaultProjectID
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	defer client.Close()

	now := time.Now().UTC()
	// ローカル検証用ユーザーの初期パスワードは全員 seedPassword にする。
	passwordHash, err := auth.HashPassword(seedPassword)
	if err != nil {
		log.Fatalf("failed to hash seed password: %v", err)
	}

	users := []UserSeed{
		{
			ID:           testUserID,
			Name:         "Test User",
			Email:        "test@example.com",
			Nickname:     "test",
			PasswordHash: passwordHash,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           "user-002",
			Name:         "User Two",
			Email:        "user2@example.com",
			Nickname:     "user2",
			PasswordHash: passwordHash,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	sessions := []SessionSeed{
		{
			ID:                            "session-test-user-001",
			UserID:                        users[0].ID,
			ProviderType:                  "local",
			InternalJWTJTI:                "seed-jti-test-user-001",
			InternalAccessTokenExpiresAt:  now.Add(15 * time.Minute),
			ProviderAccessTokenEncrypted:  "seed-encrypted-provider-access-token",
			ProviderAccessTokenExpiresAt:  now.Add(15 * time.Minute),
			ProviderRefreshTokenEncrypted: "seed-encrypted-provider-refresh-token",
			ProviderRefreshTokenExpiresAt: now.Add(24 * time.Hour),
			CreatedAt:                     now,
			UpdatedAt:                     now,
		},
	}

	for _, user := range users {
		if _, err := client.Collection("users").Doc(user.ID).Set(ctx, user); err != nil {
			log.Fatalf("failed to seed user %s: %v", user.ID, err)
		}
		log.Printf("seeded user: %s", user.ID)
	}

	for _, session := range sessions {
		if _, err := client.Collection("sessions").Doc(session.ID).Set(ctx, session); err != nil {
			log.Fatalf("failed to seed session %s: %v", session.ID, err)
		}
		log.Printf("seeded session: %s", session.ID)
	}

	log.Println("seed completed")
}

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
	// ローカル検証用ユーザーの初期パスワードは全員 password にする。
	passwordHash, err := auth.HashPassword("password")
	if err != nil {
		log.Fatalf("failed to hash seed password: %v", err)
	}

	users := []UserSeed{
		{
			ID:           "test-user",
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
			ID:                    "session-test-user-001",
			UserID:                "test-user",
			AccessTokenEncrypted:  "seed-encrypted-access-token",
			AccessTokenExpiresAt:  now.Add(15 * time.Minute),
			RefreshTokenEncrypted: "seed-encrypted-refresh-token",
			RefreshTokenExpiresAt: now.Add(24 * time.Hour),
			InternalJWTJTI:        "seed-jti-test-user-001",
			ProviderType:          "local",
			CreatedAt:             now,
			UpdatedAt:             now,
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

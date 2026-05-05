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

type PostSeed struct {
	ID           string    `firestore:"id"`
	UserID       string    `firestore:"user_id"`
	Content      string    `firestore:"content"`
	LikeCount    int       `firestore:"like_count"`
	CommentCount int       `firestore:"comment_count"`
	CreatedAt    time.Time `firestore:"created_at"`
	UpdatedAt    time.Time `firestore:"updated_at"`
}

type FollowSeed struct {
	ID         string    `firestore:"id"`
	UserIDFrom string    `firestore:"user_id_from"`
	UserIDTo   string    `firestore:"user_id_to"`
	CreatedAt  time.Time `firestore:"created_at"`
	UpdatedAt  time.Time `firestore:"updated_at"`
}

type TimelineSeed struct {
	ID                  string    `firestore:"id"`
	PostID              string    `firestore:"post_id"`
	PostUserID          string    `firestore:"post_user_id"`
	Content             string    `firestore:"content"`
	PostCreatedAt       time.Time `firestore:"post_created_at"`
	TimelineOwnerUserID string    `firestore:"timeline_owner_user_id"`
	CreatedAt           time.Time `firestore:"created_at"`
	UpdatedAt           time.Time `firestore:"updated_at"`
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
		{
			ID:           "user-003",
			Name:         "User Three",
			Email:        "user3@example.com",
			Nickname:     "user3",
			PasswordHash: passwordHash,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           "user-004",
			Name:         "User Four",
			Email:        "user4@example.com",
			Nickname:     "user4",
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

	posts := []PostSeed{
		{
			ID:           "post-test-user-001",
			UserID:       users[0].ID,
			Content:      "test-user からの最初の投稿です。",
			LikeCount:    1,
			CommentCount: 0,
			CreatedAt:    now.Add(-2 * time.Hour),
			UpdatedAt:    now.Add(-2 * time.Hour),
		},
		{
			ID:           "post-user-002-001",
			UserID:       users[1].ID,
			Content:      "user-002 の投稿です。タイムライン表示確認用です。",
			LikeCount:    2,
			CommentCount: 1,
			CreatedAt:    now.Add(-1 * time.Hour),
			UpdatedAt:    now.Add(-1 * time.Hour),
		},
		{
			ID:           "post-test-user-002",
			UserID:       users[0].ID,
			Content:      "タイムライン API を試すための 2 件目の投稿です。",
			LikeCount:    0,
			CommentCount: 0,
			CreatedAt:    now.Add(-30 * time.Minute),
			UpdatedAt:    now.Add(-30 * time.Minute),
		},
		{
			ID:           "post-user-002-002",
			UserID:       users[1].ID,
			Content:      "event driven な流れを意識して投稿を増やしています。",
			LikeCount:    3,
			CommentCount: 0,
			CreatedAt:    now.Add(-20 * time.Minute),
			UpdatedAt:    now.Add(-20 * time.Minute),
		},
		{
			ID:           "post-user-003-001",
			UserID:       users[2].ID,
			Content:      "user-003 からのはじめての投稿です。",
			LikeCount:    4,
			CommentCount: 2,
			CreatedAt:    now.Add(-90 * time.Minute),
			UpdatedAt:    now.Add(-90 * time.Minute),
		},
		{
			ID:           "post-user-003-002",
			UserID:       users[2].ID,
			Content:      "フォロー中ユーザーが複数いるケースの確認用です。",
			LikeCount:    1,
			CommentCount: 0,
			CreatedAt:    now.Add(-10 * time.Minute),
			UpdatedAt:    now.Add(-10 * time.Minute),
		},
		{
			ID:           "post-user-004-001",
			UserID:       users[3].ID,
			Content:      "この投稿は test-user のタイムラインには出ない想定です。",
			LikeCount:    0,
			CommentCount: 0,
			CreatedAt:    now.Add(-15 * time.Minute),
			UpdatedAt:    now.Add(-15 * time.Minute),
		},
	}

	follows := []FollowSeed{
		{
			ID:         "follow-test-user-user-002",
			UserIDFrom: users[0].ID,
			UserIDTo:   users[1].ID,
			CreatedAt:  now.Add(-3 * time.Hour),
			UpdatedAt:  now.Add(-3 * time.Hour),
		},
		{
			ID:         "follow-test-user-user-003",
			UserIDFrom: users[0].ID,
			UserIDTo:   users[2].ID,
			CreatedAt:  now.Add(-150 * time.Minute),
			UpdatedAt:  now.Add(-150 * time.Minute),
		},
		{
			ID:         "follow-user-002-user-004",
			UserIDFrom: users[1].ID,
			UserIDTo:   users[3].ID,
			CreatedAt:  now.Add(-80 * time.Minute),
			UpdatedAt:  now.Add(-80 * time.Minute),
		},
	}

	timelines := []TimelineSeed{
		{
			ID:                  "timeline-test-user-post-user-002-001",
			PostID:              posts[1].ID,
			PostUserID:          posts[1].UserID,
			Content:             posts[1].Content,
			PostCreatedAt:       posts[1].CreatedAt,
			TimelineOwnerUserID: users[0].ID,
			CreatedAt:           now.Add(-59 * time.Minute),
			UpdatedAt:           now.Add(-59 * time.Minute),
		},
		{
			ID:                  "timeline-test-user-post-test-user-001",
			PostID:              posts[0].ID,
			PostUserID:          posts[0].UserID,
			Content:             posts[0].Content,
			PostCreatedAt:       posts[0].CreatedAt,
			TimelineOwnerUserID: users[0].ID,
			CreatedAt:           now.Add(-119 * time.Minute),
			UpdatedAt:           now.Add(-119 * time.Minute),
		},
		{
			ID:                  "timeline-test-user-post-user-003-001",
			PostID:              posts[4].ID,
			PostUserID:          posts[4].UserID,
			Content:             posts[4].Content,
			PostCreatedAt:       posts[4].CreatedAt,
			TimelineOwnerUserID: users[0].ID,
			CreatedAt:           now.Add(-89 * time.Minute),
			UpdatedAt:           now.Add(-89 * time.Minute),
		},
		{
			ID:                  "timeline-test-user-post-test-user-002",
			PostID:              posts[2].ID,
			PostUserID:          posts[2].UserID,
			Content:             posts[2].Content,
			PostCreatedAt:       posts[2].CreatedAt,
			TimelineOwnerUserID: users[0].ID,
			CreatedAt:           now.Add(-29 * time.Minute),
			UpdatedAt:           now.Add(-29 * time.Minute),
		},
		{
			ID:                  "timeline-test-user-post-user-002-002",
			PostID:              posts[3].ID,
			PostUserID:          posts[3].UserID,
			Content:             posts[3].Content,
			PostCreatedAt:       posts[3].CreatedAt,
			TimelineOwnerUserID: users[0].ID,
			CreatedAt:           now.Add(-19 * time.Minute),
			UpdatedAt:           now.Add(-19 * time.Minute),
		},
		{
			ID:                  "timeline-test-user-post-user-003-002",
			PostID:              posts[5].ID,
			PostUserID:          posts[5].UserID,
			Content:             posts[5].Content,
			PostCreatedAt:       posts[5].CreatedAt,
			TimelineOwnerUserID: users[0].ID,
			CreatedAt:           now.Add(-9 * time.Minute),
			UpdatedAt:           now.Add(-9 * time.Minute),
		},
		{
			ID:                  "timeline-user-002-post-user-004-001",
			PostID:              posts[6].ID,
			PostUserID:          posts[6].UserID,
			Content:             posts[6].Content,
			PostCreatedAt:       posts[6].CreatedAt,
			TimelineOwnerUserID: users[1].ID,
			CreatedAt:           now.Add(-14 * time.Minute),
			UpdatedAt:           now.Add(-14 * time.Minute),
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

	for _, post := range posts {
		if _, err := client.Collection("posts").Doc(post.ID).Set(ctx, post); err != nil {
			log.Fatalf("failed to seed post %s: %v", post.ID, err)
		}
		log.Printf("seeded post: %s", post.ID)
	}

	for _, follow := range follows {
		if _, err := client.Collection("follows").Doc(follow.ID).Set(ctx, follow); err != nil {
			log.Fatalf("failed to seed follow %s: %v", follow.ID, err)
		}
		log.Printf("seeded follow: %s", follow.ID)
	}

	for _, timeline := range timelines {
		if _, err := client.Collection("timelines").Doc(timeline.ID).Set(ctx, timeline); err != nil {
			log.Fatalf("failed to seed timeline %s: %v", timeline.ID, err)
		}
		log.Printf("seeded timeline: %s", timeline.ID)
	}

	log.Println("seed completed")
}

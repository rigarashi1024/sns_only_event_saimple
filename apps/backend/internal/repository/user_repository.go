package repository

import (
	"context"
	"errors"

	gofirestore "cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID       string `firestore:"id"`
	Name     string `firestore:"name"`
	Email    string `firestore:"email"`
	Nickname string `firestore:"nickname"`
}

type UserRepository struct {
	client *gofirestore.Client
}

func NewUserRepository(client *gofirestore.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) FindByLoginID(ctx context.Context, loginID string) (*User, error) {
	doc, err := r.client.Collection("users").Doc(loginID).Get(ctx)
	if err == nil {
		var user User
		if decodeErr := doc.DataTo(&user); decodeErr != nil {
			return nil, decodeErr
		}
		if user.ID == "" {
			user.ID = doc.Ref.ID
		}
		return &user, nil
	}
	if status.Code(err) != codes.NotFound {
		return nil, err
	}

	iter := r.client.Collection("users").Where("email", "==", loginID).Limit(1).Documents(ctx)
	defer iter.Stop()

	snap, iterErr := iter.Next()
	if iterErr != nil {
		if errors.Is(iterErr, iterator.Done) {
			return nil, ErrUserNotFound
		}
		return nil, iterErr
	}

	var user User
	if decodeErr := snap.DataTo(&user); decodeErr != nil {
		return nil, decodeErr
	}
	if user.ID == "" {
		user.ID = snap.Ref.ID
	}
	return &user, nil
}

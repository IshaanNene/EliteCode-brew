package db

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/IshaanNene/EliteCode-brew/internal/auth"
)

type Store struct {
	Client *firestore.Client
	Ctx    context.Context
}

// NewStore initializes Store with FirebaseDB
func NewStore(db *auth.FirebaseDB) *Store {
	return &Store{
		Client: db.Client(),
		Ctx:    context.Background(),
	}
}

func (s *Store) SaveUser(user *auth.User) error {
	user.ID = s.Client.Collection("users").NewDoc().ID
	user.JoinedAt = time.Now()
	_, err := s.Client.Collection("users").Doc(user.ID).Set(s.Ctx, user)
	return err
}

func (s *Store) GetUserByUsername(username string) (*auth.User, error) {
	iter := s.Client.Collection("users").Where("username", "==", username).Documents(s.Ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		return nil, err
	}

	var user auth.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

package auth

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID       string    `firestore:"id" json:"id"`
	Name     string    `firestore:"name" json:"name"`
	Email    string    `firestore:"email" json:"email"`
	Username string    `firestore:"username" json:"username"`
	Password string    `firestore:"password,omitempty" json:"-"`
	JoinedAt time.Time `firestore:"joinedAt" json:"joinedAt"`
}

func GenerateUUID() string {
	return uuid.New().String()
}

func CurrentUserID() (string, error) {
	s, err := LoadSession()
	if err != nil {
		return "", err
	}
	return s.UserID, nil
}

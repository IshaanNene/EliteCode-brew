package auth

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4"
	"fmt"
	"google.golang.org/api/option"
)

type FirebaseDB struct {
	client *firestore.Client
	ctx    context.Context
}

// NewFirebaseDB initializes Firebase connection
func NewFirebaseDB(credentialsPath string) (*FirebaseDB, error) {
	ctx := context.Background()

	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore client: %v", err)
	}

	return &FirebaseDB{
		client: client,
		ctx:    ctx,
	}, nil
}

// Client exposes the Firestore client
func (db *FirebaseDB) Client() *firestore.Client {
	return db.client
}

// Close closes the Firestore client connection
func (db *FirebaseDB) Close() error {
	return db.client.Close()
}

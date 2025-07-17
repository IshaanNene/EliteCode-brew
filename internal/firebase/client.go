package firebase

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/storage"
)

type Client struct {
	App       *firebase.App
	Auth      *auth.Client
	Firestore *firestore.Client
	Storage   *storage.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	app, err := InitializeApp(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing auth client: %v", err)
	}

	firestore, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore client: %v", err)
	}

	storage, err := app.Storage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing storage client: %v", err)
	}

	return &Client{
		App:       app,
		Auth:      auth,
		Firestore: firestore,
		Storage:   storage,
	}, nil
}

func (c *Client) Close() error {
	var errs []error

	if err := c.Firestore.Close(); err != nil {
		errs = append(errs, fmt.Errorf("error closing firestore client: %v", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing clients: %v", errs)
	}

	return nil
}

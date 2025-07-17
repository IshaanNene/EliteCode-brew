package firebase

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

// Client holds all Firebase service clients
type Client struct {
	App       *firebase.App
	Auth      *auth.Client
	Firestore *firestore.Client
	Storage   *storage.Client
}

// NewClient initializes all Firebase services and returns a Client
func NewClient(ctx context.Context) (*Client, error) {
	app, err := InitializeApp(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	// Initialize Auth client
	auth, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing auth client: %v", err)
	}

	// Initialize Firestore client
	firestore, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore client: %v", err)
	}

	// Initialize Storage client
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

// Close closes all client connections
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
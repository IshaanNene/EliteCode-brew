package auth

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"time"
)

func SignupUser(name, email, username, password string) error {
	ctx := context.Background()

	client, err := App.Auth(ctx)
	if err != nil {
		return err
	}

	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(name)

	user, err := client.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	// Add profile to Firestore
	firestoreClient, err := firestore.NewClient(ctx, App.ProjectID)
	if err != nil {
		return err
	}
	defer firestoreClient.Close()

	userDoc := firestoreClient.Collection("users").Doc(user.UID)
	_, err = userDoc.Set(ctx, map[string]interface{}{
		"name":     name,
		"email":    email,
		"username": username,
		"joinedAt": time.Now(),
	})
	if err != nil {
		return err
	}

	fmt.Println("User created successfully! type elitecode login or elitecode help")
	return nil
}

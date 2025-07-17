package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"firebase.google.com/go/v4/auth"
)

const (
	authEmulatorHost = "FIREBASE_AUTH_EMULATOR_HOST"
	defaultHost      = "localhost:9099"
)

// SignInWithEmailPassword signs in a user with email and password using the Firebase Auth REST API
func SignInWithEmailPassword(ctx context.Context, email, password string) (*auth.UserRecord, error) {
	// Check if we're using the emulator
	emulatorHost := os.Getenv(authEmulatorHost)
	if emulatorHost == "" {
		emulatorHost = defaultHost
	}

	// Prepare the sign-in request
	reqBody := map[string]string{
		"email":    email,
		"password": password,
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	// Make the request to the emulator
	url := fmt.Sprintf("http://%s/identitytoolkit.googleapis.com/v1/accounts:signInWithPassword", emulatorHost)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(reqJSON)))
	if err != nil {
		return nil, fmt.Errorf("error making auth request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed: %s", resp.Status)
	}

	var result struct {
		LocalID string `json:"localId"`
		Email   string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Create a UserRecord from the response
	return &auth.UserRecord{
		UserInfo: &auth.UserInfo{
			UID:   result.LocalID,
			Email: result.Email,
		},
	}, nil
}

package firebase

import (
    "context"
    "firebase.google.com/go/v4"
    "google.golang.org/api/option"
    "log"
    "os"
    "path/filepath"
)

func InitFirebase() *firebase.App {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatalf("failed to get user home dir: %v", err)
    }
    credPath := filepath.Join(homeDir, ".elitecode", "firebase-service-account.json")
    opt := option.WithCredentialsFile(credPath)

    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        log.Fatalf("error initializing firebase: %v\n", err)
    }

    return app
}

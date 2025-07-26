package firebase

import (
    "context"
    "firebase.google.com/go/v4"
    "google.golang.org/api/option"
    "log"
)

func InitFirebase() *firebase.App {
    opt := option.WithCredentialsFile("firebase-service-account.json")
    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        log.Fatalf("error initializing firebase: %v\n", err)
    }
    return app
}

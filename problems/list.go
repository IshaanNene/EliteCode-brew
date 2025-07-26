// Fixed version of problems/list.go
package problems

import (
	"fmt"
	"log"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"context"
)

func ListProblems() {
	ctx := context.Background()
	opts := option.WithCredentialsFile("firebase-service-account.json")
	client, err := firestore.NewClient(ctx, "elitecode", opts)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	docs, err := client.Collection("problems").Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Failed to get problems: %v", err)
	}

	fmt.Println("Available Problems:")
	for _, doc := range docs {
		data := doc.Data()
		fmt.Printf("ID: %s | Title: %s | Difficulty: %s\n", doc.Ref.ID, data["title"], data["difficulty"])
	}
}

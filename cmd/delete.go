package cmd

import (
	"cloud.google.com/go/firestore"
	"github.com/IshaanNene/EliteCode-brew/internal/auth"
	"log"
	"time"

func main() {
	// Initialize Firebase connection
	db, err := auth.NewFirebaseDB("config/serviceAccountKey.json")
	if err != nil {
		log.Fatal("Failed to initialize Firebase:", err)
	}
	defer db.Close()

	// Example 1: Create a user
	user := auth.User{
		ID:       "user123",
		Name:     "John Doe",
		Email:    "john@example.com",
		Username: "johndoe",
	}

	if err := db.CreateUser(user); err != nil {
		log.Printf("Error creating user: %v", err)
	}

	// Example 2: Get user by ID
	retrievedUser, err := db.GetUser("user123")
	if err != nil {
		log.Printf("Error getting user: %v", err)
	} else {
		log.Printf("Retrieved user: %+v", retrievedUser)
	}

	// Example 3: Get user by username
	userByUsername, err := db.GetUserByUsername("johndoe")
	if err != nil {
		log.Printf("Error getting user by username: %v", err)
	} else {
		log.Printf("User found: %+v", userByUsername)
	}

	// Example 4: Create a project
	project := auth.Project{
		ID:          "proj123",
		Name:        "My Awesome Project",
		Description: "A cool project built with Go",
		Owner:       "user123",
		Language:    "go",
	}

	if err := db.CreateProject(project); err != nil {
		log.Printf("Error creating project: %v", err)
	}

	// Example 5: Get projects by owner
	projects, err := db.GetProjectsByOwner("user123")
	if err != nil {
		log.Printf("Error getting projects: %v", err)
	} else {
		log.Printf("Found %d projects", len(projects))
		for _, p := range projects {
			log.Printf("Project: %s - %s", p.Name, p.Description)
		}
	}

	// Example 6: Update user
	updates := map[string]interface{}{
		"name":  "John Smith",
		"email": "johnsmith@example.com",
	}

	if err := db.UpdateUser("user123", updates); err != nil {
		log.Printf("Error updating user: %v", err)
	}

	// Example 7: Write custom data
	customData := map[string]interface{}{
		"setting":   "dark_mode",
		"value":     true,
		"updatedAt": time.Now(),
	}

	if err := db.WriteCustomData("settings", "user123", customData); err != nil {
		log.Printf("Error writing custom data: %v", err)
	}

	// Example 8: Read custom data
	settingsData, err := db.ReadCustomData("settings", "user123")
	if err != nil {
		log.Printf("Error reading custom data: %v", err)
	} else {
		log.Printf("Settings: %+v", settingsData)
	}

	// Example 9: Batch operations
	operations := []func(*firestore.WriteBatch){
		func(batch *firestore.WriteBatch) {
			batch.Set(db.Client().Collection("logs").NewDoc(), map[string]interface{}{
				"action":    "user_login",
				"userID":    "user123",
				"timestamp": time.Now(),
			})
		},
		func(batch *firestore.WriteBatch) {
			batch.Set(db.Client().Collection("logs").NewDoc(), map[string]interface{}{
				"action":    "project_created",
				"userID":    "user123",
				"projectID": "proj123",
				"timestamp": time.Now(),
			})
		},
	}

	if err := db.BatchWrite(operations); err != nil {
		log.Printf("Error in batch write: %v", err)
	}

	allUsers, err := db.GetAllUsers()
	if err != nil {
		log.Printf("Error getting all users: %v", err)
	} else {
		log.Printf("Total users: %d", len(allUsers))
	}
}

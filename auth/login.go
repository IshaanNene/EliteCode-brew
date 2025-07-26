package auth

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "strings"

    "elitecode/firebase"
)

func Login() {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter email: ")
    email, _ := reader.ReadString('\n')
    fmt.Print("Enter password: ")
    password, _ := reader.ReadString('\n')

    email = strings.TrimSpace(email)
    password = strings.TrimSpace(password)

    app := firebase.InitFirebase()
    client, err := app.Auth(context.Background())
    if err != nil {
        fmt.Println("Error initializing auth client:", err)
        return
    }

    // Lookup user by email
    user, err := client.GetUserByEmail(context.Background(), email)
    if err != nil {
        fmt.Println("User not found or error occurred.")
        return
    }

    // ⚠️ No password verification is possible with Admin SDK
    fmt.Printf("Welcome back, %s (%s)\n", user.DisplayName, user.Email)

    // Optional: Save login info locally for session (e.g., write to .elitecode/session)
}

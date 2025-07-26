package auth

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "strings"

    "elitecode/firebase"
    "firebase.google.com/go/v4/auth"
)

func Signup() {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter name: ")
    name, _ := reader.ReadString('\n')
    fmt.Print("Enter email: ")
    email, _ := reader.ReadString('\n')
    fmt.Print("Enter password: ")
    password, _ := reader.ReadString('\n')

    name = strings.TrimSpace(name)
    email = strings.TrimSpace(email)
    password = strings.TrimSpace(password)

    app := firebase.InitFirebase()
    client, err := app.Auth(context.Background())
    if err != nil {
        fmt.Println("Error initializing auth client:", err)
        return
    }

    params := (&auth.UserToCreate{}).
        Email(email).
        Password(password).
        DisplayName(name)

    u, err := client.CreateUser(context.Background(), params)
    if err != nil {
        fmt.Println("Error creating user:", err)
        return
    }

    fmt.Printf("Successfully created user: %s\n", u.UID)
}

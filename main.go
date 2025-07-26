package main

import (
    "fmt"
    "os"

    "github.com/IshaanNene/EliteCode-brew/auth"
    "github.com/spf13/cobra"
)

func main() {
    var rootCmd = &cobra.Command{Use: "elitecode"}

    var signupCmd = &cobra.Command{
        Use:   "signup",
        Short: "Sign up for EliteCode",
        Run: func(cmd *cobra.Command, args []string) {
            auth.Signup()
        },
    }

    var loginCmd = &cobra.Command{
        Use:   "login",
        Short: "Login to EliteCode",
        Run: func(cmd *cobra.Command, args []string) {
            auth.Login()
        },
    }

    rootCmd.AddCommand(signupCmd)
    rootCmd.AddCommand(loginCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

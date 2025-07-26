package main

import (
    "fmt"
    "os"
    "github.com/IshaanNene/EliteCode-brew/auth"
    "github.com/IshaanNene/EliteCode-brew/problems"
    "github.com/spf13/cobra"
    "github.com/IshaanNene/EliteCode-brew/cmd"
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

    var problemsCmd = &cobra.Command{
        Use:   "problems",
        Short: "List available problems",
        Run: func(cmd *cobra.Command, args []string) {
            problems.ListProblems()
        },
    }
    var setProblemCmd = &cobra.Command{
        Use:   "set_problem",
        Short: "Select and set a problem from the list",
        Run: func(cmd *cobra.Command, args []string) {
            problems.SetProblem()
        },
    }


    rootCmd.AddCommand(signupCmd)
    rootCmd.AddCommand(loginCmd)
    rootCmd.AddCommand(problemsCmd)
    rootCmd.AddCommand(setProblemCmd)
    rootCmd.AddCommand(runCmd)
    rootCmd.AddCommand(submitCmd)
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

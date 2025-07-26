package main

import (
    "fmt"
    "os"

    "github.com/IshaanNene/EliteCode-brew/auth"
    cmdpkg "github.com/IshaanNene/EliteCode-brew/cmd"
    "github.com/IshaanNene/EliteCode-brew/problems"
    "github.com/spf13/cobra"
)

func main() {
    var rootCmd = &cobra.Command{Use: "elitecode"}

    var signupCmd = &cobra.Command{
        Use:   "signup",
        Short: "Sign up for EliteCode",
        Run: func(_ *cobra.Command, _ []string) {
            auth.Signup()
        },
    }

    var loginCmd = &cobra.Command{
        Use:   "login",
        Short: "Login to EliteCode",
        Run: func(_ *cobra.Command, _ []string) {
            auth.Login()
        },
    }

    var problemsCmd = &cobra.Command{
        Use:   "problems",
        Short: "List available problems",
        Run: func(_ *cobra.Command, _ []string) {
            problems.ListProblems()
        },
    }

    var setProblemCmd = &cobra.Command{
        Use:   "set_problem",
        Short: "Select and set a problem from the list",
        Run: func(_ *cobra.Command, _ []string) {
            problems.SetProblem()
        },
    }

    var runCmd = &cobra.Command{
        Use:   "run",
        Short: "Run the program",
        Run: func(_ *cobra.Command, _ []string) {
            problemID := problems.GetProblemID()
            lang := problems.GetSelectedLanguage()
            if problemID == "" || lang == "" {
                fmt.Println("Please select a problem and language first using `elitecode set_problem`")
                return
            }
            cmdpkg.Run()
        },
    }

    rootCmd.AddCommand(signupCmd, loginCmd, problemsCmd, setProblemCmd, runCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

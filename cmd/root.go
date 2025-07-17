package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var (
	cfgFile         string
	verbose         bool
	firebaseClient  *firebase.App
	firestoreClient *firestore.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "elitecode",
	Short: "A comprehensive coding practice platform",
	Long: `Elitecode CLI is a comprehensive coding practice platform that helps you:
- Practice coding problems
- Track your progress
- Improve your skills
- Prepare for interviews`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip Firebase initialization for certain commands
		if cmd.Name() == "help" || cmd.Name() == "version" || cmd.Name() == "completion" {
			return nil
		}

		// Initialize Firebase
		ctx := cmd.Context()

		// Get config directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting user home directory: %v", err)
		}

		configDir := filepath.Join(homeDir, ".elitecode")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("error creating config directory: %v", err)
		}

		// Initialize Firebase
		opt := option.WithCredentialsJSON([]byte(firebaseConfig))
		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil {
			return fmt.Errorf("error initializing Firebase app: %v", err)
		}
		firebaseClient = app

		// Initialize Firestore
		client, err := app.Firestore(ctx)
		if err != nil {
			return fmt.Errorf("error initializing Firestore client: %v", err)
		}
		firestoreClient = client

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		// Close Firestore client
		if firestoreClient != nil {
			if err := firestoreClient.Close(); err != nil {
				return fmt.Errorf("error closing Firestore client: %v", err)
			}
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("Error: %v", err))
		os.Exit(1)
	}
}

// ExecuteContext executes the root command with context
func ExecuteContext(ctx context.Context) {
	rootCmd.SetContext(ctx)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("Error: %v", err))
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.elitecode/config.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")

	// Hide help command
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".elitecode" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".elitecode"))
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// Firebase configuration
const firebaseConfig = `{
	"apiKey": "AIzaSyDhoR2EA3uuWqWSkiiE5G1vRkeqoGHBdro",
	"authDomain": "elitecode-47058.firebaseapp.com",
	"projectId": "elitecode-47058",
	"storageBucket": "elitecode-47058.firebasestorage.app",
	"messagingSenderId": "605884776680",
	"appId": "1:605884776680:web:f89db0d49b021139dec95d",
	"measurementId": "G-4XKNK2CH3E"
}`

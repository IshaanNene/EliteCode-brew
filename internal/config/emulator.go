package config

import (
	"fmt"
	"os"
)

const (
	FirebaseAuthEmulatorHost    = "FIREBASE_AUTH_EMULATOR_HOST"
	FirebaseFirestoreEmulator   = "FIRESTORE_EMULATOR_HOST"
	FirebaseStorageEmulatorHost = "FIREBASE_STORAGE_EMULATOR_HOST"

	DefaultAuthPort      = 9099
	DefaultFirestorePort = 8080
	DefaultStoragePort   = 9199
)

func SetupEmulators() error {
	if os.Getenv(FirebaseAuthEmulatorHost) == "" {
		if err := os.Setenv(FirebaseAuthEmulatorHost, fmt.Sprintf("localhost:%d", DefaultAuthPort)); err != nil {
			return fmt.Errorf("error setting auth emulator host: %v", err)
		}
	}

	if os.Getenv(FirebaseFirestoreEmulator) == "" {
		if err := os.Setenv(FirebaseFirestoreEmulator, fmt.Sprintf("localhost:%d", DefaultFirestorePort)); err != nil {
			return fmt.Errorf("error setting firestore emulator host: %v", err)
		}
	}

	if os.Getenv(FirebaseStorageEmulatorHost) == "" {
		if err := os.Setenv(FirebaseStorageEmulatorHost, fmt.Sprintf("localhost:%d", DefaultStoragePort)); err != nil {
			return fmt.Errorf("error setting storage emulator host: %v", err)
		}
	}

	return nil
}

func IsEmulatorEnabled() bool {
	return os.Getenv(FirebaseAuthEmulatorHost) != "" ||
		os.Getenv(FirebaseFirestoreEmulator) != "" ||
		os.Getenv(FirebaseStorageEmulatorHost) != ""
}

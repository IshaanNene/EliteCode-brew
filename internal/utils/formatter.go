package utils

import (
	"fmt"
)

func PrintSuccess(msg string) {
	fmt.Printf("✅ %s\n", msg)
}

func PrintError(err error) {
	fmt.Printf("❌ %s\n", err)
}

func PrintWarning(msg string) {
	fmt.Printf("⚠️  %s\n", msg)
}

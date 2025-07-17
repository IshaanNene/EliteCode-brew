package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/IshaanNene/EliteCode-brew/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	cmd.ExecuteContext(ctx)
}

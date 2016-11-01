package command

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func waitSignal(ctx context.Context, cancel context.CancelFunc, done chan error) error {
	sig := make(chan os.Signal, 4)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-sig:
		cancel()
		return fmt.Errorf("signal: %v", s)
	case exit := <-done:
		return exit
	case <-ctx.Done():
		return ctx.Err()
	}
}

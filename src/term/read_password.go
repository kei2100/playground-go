package term

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

// ReadPasswordContext reads a password from stdin using golang.org/x/term.ReadPassword with Context
func ReadPasswordContext(ctx context.Context) (string, error) {
	inputCh := make(chan string, 1)
	errCh := make(chan error, 1)
	stdin := syscall.Stdin
	state, err := term.GetState(stdin)
	if err != nil {
		return "", fmt.Errorf("term: get term state stdin: %w", err)
	}
	go func() {
		b, err := term.ReadPassword(stdin)
		if err != nil {
			errCh <- fmt.Errorf("term: read password from stdin: %w", err)
			return
		}
		inputCh <- string(b)
	}()
	select {
	case <-ctx.Done():
		if err := term.Restore(stdin, state); err != nil {
			fmt.Fprintf(os.Stderr, "term: failed to restore term state: %+v", err)
		}
		return "", ctx.Err()
	case err := <-errCh:
		return "", err
	case input := <-inputCh:
		return input, nil
	}
}

package wait

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// WGroup waits wg.Wait().
// When timeout exceeded while waiting for wg.Wait(), return an error.
// NOTICE: This method will "leak" a goroutine if wg.Done() does not complete.
func WGroup(wg *sync.WaitGroup, timeout time.Duration) error {
	if timeout == 0 {
		return errors.New("wait: timeout must be greater than zero")
	}
	if wg == nil {
		return errors.New("wait: arg WGroup is nil")
	}
	ctx, can := context.WithTimeout(context.Background(), timeout)
	defer can()
	return WGroupContext(ctx, wg)
}

// WGroupContext waits wg.Wait() with context.
func WGroupContext(ctx context.Context, wg *sync.WaitGroup) error {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ReceiveStructContext(ctx, ch)
}

// ReceiveStruct waits receive from the channel.
// When timeout exceeded while waiting for receive, return an error.
func ReceiveStruct(ch <-chan struct{}, timeout time.Duration) error {
	if timeout == 0 {
		return errors.New("wait: timeout must be greater than zero")
	}
	ctx, can := context.WithTimeout(context.Background(), timeout)
	defer can()
	return ReceiveStructContext(ctx, ch)
}

// ReceiveStructContext waits receive from the channel with context.
func ReceiveStructContext(ctx context.Context, ch <-chan struct{}) error {
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait: %v", ctx.Err())
	}
}

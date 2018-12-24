package pipeline2

import (
	"context"
	"testing"
	"time"
)

func generator(ctx context.Context, values ...interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for _, v := range values {
			select {
			case <-ctx.Done():
				return
			case ch <- v:
			}
		}

	}()
	return ch
}

func assertRecvSeq(t *testing.T, timeout time.Duration, ch <-chan interface{}, wants ...interface{}) {
	t.Helper()
	if timeout == 0 {
		timeout = time.Second
	}
	done := time.After(timeout)
	for i, want := range wants {
		select {
		case <-done:
			t.Errorf("timeout while waiting for %v", wants[i:])
			return
		case got, ok := <-ch:
			if !ok {
				t.Errorf("ch closed. want %v", want)
			} else if got != want {
				t.Errorf("idx %v recv %v, want %v", i, got, want)
			}
		}
	}
}

func TestGenerator(t *testing.T) {
	ctx, can := context.WithTimeout(context.Background(), time.Second)
	defer can()
	ch := generator(ctx, 1, 2, 3)
	assertRecvSeq(t, time.Second, ch, 1, 2, 3)
}

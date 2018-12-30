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

func repeat(ctx context.Context, values ...interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			for _, v := range values {
				select {
				case <-ctx.Done():
					return
				case ch <- v:
				}
			}
		}
	}()
	return ch
}

func take(ctx context.Context, stream <-chan interface{}, num int) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		done := ctx.Done()
		for i := 0; i < num; {
			select {
			case <-done:
				return
			case v := <-stream:
				select {
				case <-done:
					return
				case ch <- v:
					i++
				}
			}
		}
	}()
	return ch
}

func TestRepeatTake(t *testing.T) {
	ctx, can := context.WithTimeout(context.Background(), time.Second)
	defer can()

	ch := take(ctx, repeat(ctx, 1, 2, 3), 4)
	cnt := 0
	for v := range ch {
		cnt += v.(int)
	}
	if g, w := cnt, 7; g != w {
		t.Errorf("cnt got %v, want %v", g, w)
	}
}

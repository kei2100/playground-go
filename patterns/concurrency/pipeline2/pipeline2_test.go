package pipeline2

import (
	"context"
	"sync"
	"testing"
	"time"
	"unicode/utf8"
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
				if ok := sendOrDone(ctx, ch, v); !ok {
					return
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
		i := 0
		for v := range recvOrDone(ctx, stream) {
			if i >= num {
				return
			}
			sendOrDone(ctx, ch, v)
			i++
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

	ch = take(ctx, take(ctx, repeat(ctx, 1), 2), 4)
	cnt = 0
	for v := range ch {
		cnt += v.(int)
	}
	if g, w := cnt, 2; g != w {
		t.Errorf("cnt got %v, want %v", g, w)
	}
}

func sendOrDone(ctx context.Context, to chan<- interface{}, v interface{}) bool {
	select {
	case <-ctx.Done():
		return false
	case to <- v:
		return true
	}
}

func recvOrDone(ctx context.Context, from <-chan interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-from:
				if !ok {
					return
				}
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

func wordCount(ctx context.Context, stringStream <-chan interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for v := range recvOrDone(ctx, stringStream) {
			n := utf8.RuneCountInString(v.(string))
			sendOrDone(ctx, ch, n)
		}
	}()
	return ch
}

func merge(ctx context.Context, streams ...<-chan interface{}) <-chan interface{} {
	wg := sync.WaitGroup{}
	merged := make(chan interface{})

	mux := func(ch <-chan interface{}) {
		defer wg.Done()
		for v := range recvOrDone(ctx, ch) {
			sendOrDone(ctx, merged, v)
		}
	}

	for _, s := range streams {
		wg.Add(1)
		go mux(s)
	}
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func TestFanOutFanIn(t *testing.T) {
	ctx, can := context.WithTimeout(context.Background(), time.Second)
	defer can()

	words := generator(ctx, "foo", "barb", "bazzz")
	// fan out
	ch1 := wordCount(ctx, words)
	ch2 := wordCount(ctx, words)
	ch3 := wordCount(ctx, words)

	// fan in
	sum := 0
	for cnt := range merge(ctx, ch1, ch2, ch3) {
		sum += cnt.(int)
	}
	if g, w := sum, 12; g != w {
		t.Errorf("sum got %v, want %v", g, w)
	}
}

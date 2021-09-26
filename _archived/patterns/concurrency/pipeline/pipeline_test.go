package pipeline

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

func assertRecv(t *testing.T, ch <-chan interface{}, want interface{}) {
	t.Helper()
	ctx, can := context.WithTimeout(context.Background(), time.Second)
	defer can()
	assertRecvContext(ctx, t, ch, want)
}

func assertRecvContext(ctx context.Context, t *testing.T, ch <-chan interface{}, want interface{}) {
	t.Helper()
	done := ctx.Done()
	select {
	case <-done:
		t.Errorf("timeout while waiting for %v", want)
		return
	case got, ok := <-ch:
		if !ok {
			t.Errorf("ch closed. want %v", want)
		} else if got != want {
			t.Errorf("recv %v, want %v", got, want)
		}
	}
}

func assertRecvSeq(t *testing.T, ch <-chan interface{}, wants ...interface{}) {
	t.Helper()
	ctx, can := context.WithTimeout(context.Background(), time.Second)
	defer can()
	assertRecvSeqContext(ctx, t, ch, wants...)
}

func assertRecvSeqContext(ctx context.Context, t *testing.T, ch <-chan interface{}, wants ...interface{}) {
	t.Helper()
	done := ctx.Done()
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
	assertRecvSeq(t, ch, 1, 2, 3)
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

	drain := func(ch <-chan interface{}) {
		defer wg.Done()
		for v := range recvOrDone(ctx, ch) {
			sendOrDone(ctx, merged, v)
		}
	}

	for _, s := range streams {
		wg.Add(1)
		go drain(s)
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

func tee(ctx context.Context, stream <-chan interface{}) (_, _ chan interface{}) {
	ch1, ch2 := make(chan interface{}), make(chan interface{})
	go func() {
		defer func() {
			close(ch1)
			close(ch2)
		}()

		for v := range recvOrDone(ctx, stream) {
			var c1, c2 = ch1, ch2
			for i := 0; i < 2; i++ {
				select {
				case <-ctx.Done():
					return
				case c1 <- v:
					c1 = nil
				case c2 <- v:
					c2 = nil
				}
			}
		}
	}()
	return ch1, ch2
}

func TestTee(t *testing.T) {
	ch := make(chan interface{})
	ctx, can := context.WithTimeout(context.Background(), time.Second)
	defer can()

	go func() {
		sendOrDone(ctx, ch, 1)
		sendOrDone(ctx, ch, 2)
		sendOrDone(ctx, ch, 3)
	}()

	tc1, tc2 := tee(ctx, ch)
	assertRecv(t, tc1, 1)
	assertRecv(t, tc2, 1)
	assertRecv(t, tc1, 2)
	assertRecv(t, tc2, 2)
	assertRecv(t, tc1, 3)
	assertRecv(t, tc2, 3)
}

func bridge(ctx context.Context, chStream <-chan <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			var in <-chan interface{}
			select {
			case <-ctx.Done():
				return
			case ch, ok := <-chStream:
				if !ok {
					return
				}
				in = ch
			}
			for v := range recvOrDone(ctx, in) {
				sendOrDone(ctx, out, v)
			}
		}
	}()
	return out
}

func TestBridge(t *testing.T) {
	ctx, can := context.WithTimeout(context.Background(), time.Second)
	defer can()

	ch1 := take(ctx, repeat(ctx, 1), 1)
	ch2 := take(ctx, repeat(ctx, 2), 2)
	chStream := make(chan (<-chan interface{}))
	go func() {
		defer close(chStream)
		select {
		case <-ctx.Done():
			return
		case chStream <- ch1:
		}
		select {
		case <-ctx.Done():
			return
		case chStream <- ch2:
		}
	}()

	out := bridge(ctx, chStream)
	assertRecvSeqContext(ctx, t, out, 1, 2, 2)
}

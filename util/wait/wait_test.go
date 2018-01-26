package wait

import (
	"sync"
	"testing"
	"time"
)

func TestWaitGroup(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			wg.Done()
		}()
		if err := WGroup(&wg, 10*time.Millisecond); err != nil {
			t.Errorf("got %v, want no error", err)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			time.Sleep(50 * time.Millisecond)
			wg.Done()
		}()
		if err := WGroup(&wg, 10*time.Millisecond); err == nil {
			t.Error("got nil, want timeout error")
		}
	})
}

func TestReceive(t *testing.T) {
	t.Parallel()

	t.Run("Struct ok", func(t *testing.T) {
		ch := make(chan struct{}, 1)
		go func() { ch <- struct{}{} }()
		if err := ReceiveStruct(ch, 10*time.Millisecond); err != nil {
			t.Errorf("got %v, want no error", err)
		}
	})

	t.Run("Struct timeout", func(t *testing.T) {
		ch := make(chan struct{}, 1)
		if err := ReceiveStruct(ch, 10*time.Millisecond); err == nil {
			t.Error("got nil, want timeout error")
		}
	})
}

package wait

import (
	"testing"
	"time"
)

func TestReceive_WithTimeout(t *testing.T) {
	t.Parallel()

	t.Run("Struct ok", func(t *testing.T) {
		ch := make(chan struct{}, 1)
		go func() { ch <- struct{}{} }()
		err := ReceiveStuct(ch, WithTimeout(100*time.Millisecond))
		if err != nil {
			t.Errorf("got %v, want no error", err)
		}
	})

	t.Run("Struct timeout", func(t *testing.T) {
		ch := make(chan struct{}, 1)
		err := ReceiveStuct(ch, WithTimeout(100*time.Millisecond))
		if err == nil {
			t.Error("got nil, want timeout error")
		}
	})
}

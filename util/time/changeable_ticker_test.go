package time

import (
	"testing"
	"time"
)

func TestChangeableTicker(t *testing.T) {
	t.Parallel()

	done := make(chan struct{})
	startAt := time.Now()

	go func() {
		ticker := NewChangeableTicker(10 * time.Millisecond)
		defer ticker.Stop()
		<-ticker.C()
		<-ticker.C()
		ticker.Change(100 * time.Millisecond)
		<-ticker.C()
		close(done)
	}()

	timeout := time.Second
	to := time.After(timeout)

	select {
	case <-to:
		t.Errorf("timeout %s exceeded while waiting for receive ticks", timeout)
	case <-done:
		break
	}

	elapsed := time.Now().Sub(startAt)
	if elapsed < 120*time.Millisecond {
		t.Errorf("unexpected elapsed %s", elapsed)
	}
}

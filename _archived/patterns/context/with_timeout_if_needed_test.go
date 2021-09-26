package context

import (
	"context"
	"testing"
	"time"
)

var noop = func() {}

func WithTimeoutIfNeeded(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if _, ok := parent.Deadline(); ok {
		return parent, noop
	}
	return context.WithTimeout(parent, timeout)
}

func TestWithTimeoutIfNeeded(t *testing.T) {
	t.Run("parent is deadline context", func(t *testing.T) {
		p, pcan := context.WithTimeout(context.Background(), time.Millisecond)
		defer pcan()

		c, can := WithTimeoutIfNeeded(p, time.Second)
		defer can()

		select {
		case <-c.Done():
			if g, w := c.Err(), context.DeadlineExceeded; g != w {
				t.Errorf("c.Err() got %v, want %v", g, w)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("unexpected timeout while waiting for c.Done()")
		}
	})
	t.Run("parent is cancel context", func(t *testing.T) {
		p, pcan := context.WithCancel(context.Background())
		defer pcan()

		c, can := WithTimeoutIfNeeded(p, time.Millisecond)
		defer can()

		select {
		case <-c.Done():
			if g, w := c.Err(), context.DeadlineExceeded; g != w {
				t.Errorf("c.Err() got %v, want %v", g, w)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("unexpected timeout while waiting for c.Done()")
		}
	})
	t.Run("parent is background context", func(t *testing.T) {
		c, can := WithTimeoutIfNeeded(context.Background(), time.Millisecond)
		defer can()

		select {
		case <-c.Done():
			if g, w := c.Err(), context.DeadlineExceeded; g != w {
				t.Errorf("c.Err() got %v, want %v", g, w)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("unexpected timeout while waiting for c.Done()")
		}
	})
}

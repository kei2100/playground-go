package context

import (
	"context"
	"testing"
)

func assertDoneWithCanceled(c context.Context, t *testing.T) {
	t.Helper()

	select {
	case <-c.Done():
		if g, w := c.Err(), context.Canceled; g != w {
			t.Errorf("Err got %v, want %v", g, w)
		}
	default:
		t.Error("Done channel got not received, want received")
	}
}

func assertNotDone(c context.Context, t *testing.T) {
	t.Helper()

	select {
	case <-c.Done():
		t.Error("Done channel got received, want not received")
	default:
		// ok
	}
}

func TestContextWithCancel(t *testing.T) {
	t.Run("duplicate call cancel", func(t *testing.T) {
		c, can := context.WithCancel(context.Background())

		can()
		assertDoneWithCanceled(c, t)
		can()
		assertDoneWithCanceled(c, t)
	})

	t.Run("when parent canceled", func(t *testing.T) {
		p, pcan := context.WithCancel(context.Background())
		c, can := context.WithCancel(p)
		defer can()

		pcan()
		assertDoneWithCanceled(p, t)
		assertDoneWithCanceled(c, t)
	})

	t.Run("when child canceled", func(t *testing.T) {
		p, pcan := context.WithCancel(context.Background())
		c, can := context.WithCancel(p)
		defer pcan()

		can()
		assertNotDone(p, t)
		assertDoneWithCanceled(c, t)
	})
}

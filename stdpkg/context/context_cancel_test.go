package context

import (
	"context"
	"testing"
)

func assertDoneWithCanceled(t *testing.T, c context.Context) {
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

func assertNotDone(t *testing.T, c context.Context) {
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
		assertDoneWithCanceled(t, c)
		can()
		assertDoneWithCanceled(t, c)
	})

	t.Run("when parent canceled", func(t *testing.T) {
		p, pcan := context.WithCancel(context.Background())
		c, can := context.WithCancel(p)
		defer can()

		pcan()
		assertDoneWithCanceled(t, p)
		assertDoneWithCanceled(t, c)
	})

	t.Run("when child canceled", func(t *testing.T) {
		p, pcan := context.WithCancel(context.Background())
		c, can := context.WithCancel(p)
		defer pcan()

		can()
		assertNotDone(t, p)
		assertDoneWithCanceled(t, c)
	})
}

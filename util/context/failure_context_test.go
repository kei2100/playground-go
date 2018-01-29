package context

import (
	"context"
	"errors"
	"testing"
	"time"
)

func assertDone(t *testing.T, c context.Context, wanterr error) {
	t.Helper()

	select {
	case <-c.Done():
		if g, w := c.Err(), wanterr; g != w {
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

func TestFailureContext_Cancel(t *testing.T) {
	t.Run("call cancel", func(t *testing.T) {
		c, can := Failure(context.Background())
		assertNotDone(t, c)
		can()
		assertDone(t, c, context.Canceled)
	})

	t.Run("duplicate call cancel", func(t *testing.T) {
		c, can := Failure(context.Background())
		can()
		assertDone(t, c, context.Canceled)
		can()
		assertDone(t, c, context.Canceled)
	})
}

func TestFailureContext_Fail(t *testing.T) {
	t.Run("call Fail", func(t *testing.T) {
		c, can := Failure(context.Background())
		defer can()

		wanterr := errors.New("want error")
		c.Fail(wanterr)
		assertDone(t, c, wanterr)
	})

	t.Run("duplicate call Fail", func(t *testing.T) {
		c, can := Failure(context.Background())
		defer can()

		want := errors.New("want error")
		c.Fail(want)
		assertDone(t, c, want)
		c.Fail(errors.New("unexpected"))
		assertDone(t, c, want)
	})
}

func TestFailureContext_Propagation(t *testing.T) {
	t.Run("parent to child cancel", func(t *testing.T) {
		p, pcan := context.WithCancel(context.Background())
		c, can := Failure(p)
		defer can()
		pcan()
		assertDone(t, p, context.Canceled)
		assertDone(t, c, context.Canceled)
	})

	t.Run("parent to child failure", func(t *testing.T) {
		p, pcan := context.WithTimeout(context.Background(), time.Millisecond)
		defer pcan()
		c, can := Failure(p)
		defer can()

		time.Sleep(10 * time.Millisecond)
		assertDone(t, p, context.DeadlineExceeded)
		assertDone(t, c, context.DeadlineExceeded)
	})

	t.Run("child to parent cancel", func(t *testing.T) {
		p, pcan := context.WithCancel(context.Background())
		defer pcan()
		c, can := Failure(p)
		can()
		assertNotDone(t, p)
		assertDone(t, c, context.Canceled)
	})

	t.Run("child to parent Fail", func(t *testing.T) {
		p, pcan := context.WithCancel(context.Background())
		defer pcan()
		c, can := Failure(p)
		defer can()

		want := errors.New("want error")
		c.Fail(want)
		assertNotDone(t, p)
		assertDone(t, c, want)
	})
}

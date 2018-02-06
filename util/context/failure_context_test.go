package context

import (
	"context"
	"errors"
	"testing"
	"time"
)

func assertDone(c context.Context, t *testing.T, wanterr error) {
	t.Helper()

	select {
	case <-c.Done():
		if g, w := c.Err(), wanterr; g != w {
			t.Errorf("Err got %v, want %v", g, w)
		}
	case <-time.After(10 * time.Millisecond):
		t.Error("Done channel got not received, want received")
	}
}

func assertNotDone(c context.Context, t *testing.T) {
	t.Helper()

	select {
	case <-c.Done():
		t.Error("Done channel got received, want not received")
	case <-time.After(10 * time.Millisecond):
		// ok
	}
}

func TestFailureContext_Cancel(t *testing.T) {
	t.Run("call cancel", func(t *testing.T) {
		c, can := Failure(context.Background())
		assertNotDone(c, t)
		can()
		assertDone(c, t, context.Canceled)
	})

	t.Run("duplicate call cancel", func(t *testing.T) {
		c, can := Failure(context.Background())
		can()
		assertDone(c, t, context.Canceled)
		can()
		assertDone(c, t, context.Canceled)
	})
}

func TestFailureContext_Fail(t *testing.T) {
	t.Run("call Fail", func(t *testing.T) {
		c, can := Failure(context.Background())
		defer can()

		wanterr := errors.New("want error")
		c.Fail(wanterr)
		assertDone(c, t, wanterr)
	})

	t.Run("duplicate call Fail", func(t *testing.T) {
		c, can := Failure(context.Background())
		defer can()

		want := errors.New("want error")
		c.Fail(want)
		assertDone(c, t, want)
		c.Fail(errors.New("unexpected"))
		assertDone(c, t, want)
	})
}

func TestFailureContext_Propagation(t *testing.T) {
	t.Run("parent to child cancel", func(t *testing.T) {
		p, pcan := context.WithCancel(context.Background())
		c, can := Failure(p)
		defer can()
		pcan()
		assertDone(p, t, context.Canceled)
		assertDone(c, t, context.Canceled)
	})

	t.Run("parent to child failure", func(t *testing.T) {
		p, pcan := context.WithTimeout(context.Background(), time.Millisecond)
		defer pcan()
		c, can := Failure(p)
		defer can()

		time.Sleep(10 * time.Millisecond)
		assertDone(p, t, context.DeadlineExceeded)
		assertDone(c, t, context.DeadlineExceeded)
	})

	t.Run("child to parent cancel", func(t *testing.T) {
		p, pcan := context.WithCancel(context.Background())
		defer pcan()
		c, can := Failure(p)
		can()
		assertNotDone(p, t)
		assertDone(c, t, context.Canceled)
	})

	t.Run("child to parent Fail", func(t *testing.T) {
		p, pcan := context.WithCancel(context.Background())
		defer pcan()
		c, can := Failure(p)
		defer can()

		want := errors.New("want error")
		c.Fail(want)
		assertNotDone(p, t)
		assertDone(c, t, want)
	})
}

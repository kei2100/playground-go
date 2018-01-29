package context

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// 任意のerrorで終了できるカスタムコンテキスト
type errorMarkableContext struct {
	context.Context

	once sync.Once
	done chan struct{}
	err  error
}

func WithErrorMarkable(p context.Context, pCancel context.CancelFunc) (*errorMarkableContext, context.CancelFunc) {
	c := &errorMarkableContext{
		Context: p,
		done:    make(chan struct{}),
	}
	can := func() {
		if pCancel != nil {
			pCancel()
			return
		}
		c.MarkAsError(context.Canceled)
	}

	if p.Done() == nil {
		return c, can
	}
	go func() {
		select {
		case <-p.Done():
			c.MarkAsError(p.Err())
		case <-c.done:
			return
		}
	}()
	return c, can
}

func (c *errorMarkableContext) MarkAsError(err error) {
	c.once.Do(func() {
		c.err = err
		close(c.done)
	})
}

func (c *errorMarkableContext) Done() <-chan struct{} {
	return c.done
}

func (c *errorMarkableContext) Err() error {
	return c.err
}

func (c *errorMarkableContext) String() string {
	// e.g.
	// - context.Background.WithErrorMarkable()
	// - context.Background.WithDeadline(2018-01-29 11:35:01.53441438 +0900 JST m=+1.001233545 [999.956341ms]).WithErrorMarkable()
	return fmt.Sprintf("%v.WithErrorMarkable()", c.Context)
}

func TestErrorMarkableContext(t *testing.T) {
	t.Run("call cancel", func(t *testing.T) {
		tests := []struct {
			subject string
			param   func() (context.Context, context.CancelFunc)
		}{
			{
				subject: "parent is background",
				param:   func() (context.Context, context.CancelFunc) { return context.Background(), nil },
			},
			{
				subject: "parent is deadline",
				param: func() (context.Context, context.CancelFunc) {
					return context.WithDeadline(context.Background(), time.Now().Add(time.Second))
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.subject, func(t *testing.T) {
				ctx, can := WithErrorMarkable(tt.param())
				select {
				case <-ctx.Done():
					t.Error("Done() chan got received, want not received")
				default:
					// ok
				}

				go can()

				select {
				case <-ctx.Done():
					if g, w := ctx.Err(), context.Canceled; g != w {
						t.Errorf("Err() got '%v', want '%v'", g, w)
					}
				case <-time.After(10 * time.Millisecond):
					t.Error("timeout exceeded while waiting for Done chan received")
				}
			})
		}
	})

	t.Run("call MarkAsError", func(t *testing.T) {
		tests := []struct {
			subject string
			param   func() (context.Context, context.CancelFunc)
		}{
			{
				subject: "parent is background",
				param:   func() (context.Context, context.CancelFunc) { return context.Background(), nil },
			},
			{
				subject: "parent is deadline",
				param: func() (context.Context, context.CancelFunc) {
					return context.WithDeadline(context.Background(), time.Now().Add(time.Second))
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.subject, func(t *testing.T) {
				ctx, can := WithErrorMarkable(tt.param())
				defer can()

				want := errors.New("want error")
				go ctx.MarkAsError(want)

				select {
				case <-ctx.Done():
					if g, w := ctx.Err(), want; g != w {
						t.Errorf("Err() got '%v', want '%v'", g, w)
					}
				case <-time.After(10 * time.Millisecond):
					t.Error("timeout exceeded while waiting for Done chan received")
				}
			})
		}
	})

	t.Run("parent context canceled", func(t *testing.T) {
		p, pcan := context.WithTimeout(context.Background(), time.Second)
		ctx, can := WithErrorMarkable(p, pcan)
		defer can()

		pcan()
		select {
		case <-ctx.Done():
			if g, w := ctx.Err(), context.Canceled; g != w {
				t.Errorf("Err() got '%v', want '%v'", g, w)
			}
		case <-time.After(10 * time.Millisecond):
			t.Error("timeout exceeded while waiting for Done chan received")
		}
	})

	t.Run("parent context failure", func(t *testing.T) {
		ctx, can := WithErrorMarkable(context.WithTimeout(context.Background(), time.Millisecond))
		defer can()

		time.Sleep(2 * time.Millisecond)
		select {
		case <-ctx.Done():
			if g, w := ctx.Err(), context.DeadlineExceeded; g != w {
				t.Errorf("Err() got '%v', want '%v'", g, w)
			}
		case <-time.After(10 * time.Millisecond):
			t.Error("timeout exceeded while waiting for Done chan received")
		}
	})
}

package context

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// failureContext is an implementation of context.Context, which can be terminated with any error.
type failureContext struct {
	p context.Context

	once sync.Once
	done chan struct{}
	err  error
}

// Failure creates the failureContext. The returned context's Done channel is closed
// when the returned cancel function is called
// or when call Fail to the context with an error
// or when the parent context's Done channel is closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func Failure(p context.Context) (*failureContext, context.CancelFunc) {
	c := &failureContext{
		p:    p,
		done: make(chan struct{}),
	}
	can := func() {
		c.Fail(context.Canceled)
	}

	if p.Done() == nil {
		return c, can
	}
	go func() {
		select {
		case <-p.Done():
			c.Fail(p.Err())
		case <-c.done:
			return
		}
	}()
	return c, can
}

// Fail terminates this context with an error.
// if Fail is called multiple times, only the first call will set the error
func (c *failureContext) Fail(err error) {
	c.once.Do(func() {
		c.err = err
		close(c.done)
	})
}

func (c *failureContext) Deadline() (deadline time.Time, ok bool) {
	return c.p.Deadline()
}

func (c *failureContext) Done() <-chan struct{} {
	return c.done
}

func (c *failureContext) Err() error {
	return c.err
}

func (c *failureContext) Value(key interface{}) interface{} {
	return c.p.Value(key)
}

func (c *failureContext) String() string {
	// e.g.
	// - context.Background.Failure()
	// - context.Background.WithDeadline(2018-01-29 11:35:01.53441438 +0900 JST m=+1.001233545 [999.956341ms]).Failure()
	return fmt.Sprintf("%v.Failure()", c.p)
}

package once

import (
	"errors"
	"sync"
	"testing"
)

type errorCloser struct {
	closed bool
}

func (c *errorCloser) Close() error {
	defer func() { c.closed = true }()
	return errors.New("always return error")
}

func TestCloser(t *testing.T) {
	t.Parallel()

	impl := new(errorCloser)
	cl := NewCloser(impl)

	wg := sync.WaitGroup{}
	wg.Add(2)

	var r1, r2 error

	go func() {
		r1 = cl.Close()
		wg.Done()
	}()
	go func() {
		r2 = cl.Close()
		wg.Done()
	}()

	wg.Wait()

	if (r1 == nil && r2 != nil) && (r1 != nil && r2 == nil) {
		t.Errorf("got r1 %v and r2 %v", r1, r2)
	}
}

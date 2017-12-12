package concurrency

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

// OnceOrError is similar to sync.Once, but perform the action until it succeeds
type OnceOrError struct {
	m    sync.Mutex
	done uint32
}

// Do is similar to sync.Once.Do, but perform the action until it succeeds
func (o *OnceOrError) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		if err := f(); err != nil {
			return err
		}
		atomic.StoreUint32(&o.done, 1)
	}
	return nil
}

func TestOnceOrError(t *testing.T) {
	var wg sync.WaitGroup
	var count, errCount int32

	once := new(OnceOrError)

	f := func(raiseErr bool) {
		defer wg.Done()
		once.Do(func() error {
			if raiseErr {
				atomic.AddInt32(&errCount, 1)
				return errors.New("")
			}
			atomic.AddInt32(&count, 1)
			return nil
		})
	}

	const N = 10
	for i := 0; i < N; i++ {
		wg.Add(1)
		go f(true)
	}
	wg.Wait()
	if g, w := errCount, int32(N); g != w {
		t.Errorf("errCount got %v, want %v", g, w)
	}

	for i := 0; i < N; i++ {
		wg.Add(1)
		go f(false)
	}
	wg.Wait()
	if g, w := count, int32(1); g != w {
		t.Errorf("count got %v, want %v", g, w)
	}
}

package sync

import (
	"sync"
	"sync/atomic"
)

// OnceOrError is similar to sync.Once, but perform the action until it succeeds
type OnceOrError struct {
	m    sync.Mutex
	done uint32
}

// DoOrError is similar to sync.Once.Do, but perform the action until it succeeds
func (o *OnceOrError) DoOrError(f func() error) error {
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

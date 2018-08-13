package sync

import (
	"sync"
)

// MultiErrGroup internally has sync.WaitGroup,
// It keeps errors that occurred when each goroutine was executed
type MultiErrGroup struct {
	wg   sync.WaitGroup
	errs []error
	mu   sync.Mutex
}

// Wait blocks until all function calls from the Go method have returned,
// then returns all errors from them.
func (g *MultiErrGroup) Wait() []error {
	g.wg.Wait()
	return g.errs
}

// Go calls the given function in a new goroutine.
func (g *MultiErrGroup) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.mu.Lock()
			defer g.mu.Unlock()
			g.errs = append(g.errs, err)
		}
	}()
}

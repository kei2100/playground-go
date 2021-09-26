package once

import "sync"

type closer interface {
	Close() error
}

// Closer provides once Close() call
type Closer struct {
	closer closer
	once   sync.Once
}

// NewCloser returns new Closer
func NewCloser(impl closer) *Closer {
	return &Closer{closer: impl}
}

// Close once call closer.Close()
func (c *Closer) Close() error {
	var err error
	c.once.Do(func() {
		err = c.closer.Close()
	})
	return err
}

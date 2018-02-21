package context

import (
	"context"
	"sync"
	"time"
)

// TimestampContext is the context in which time stamps can be recorded
type TimestampContext interface {
	context.Context
	// Timestamp records the timestamp to this context with the given description
	Timestamp(description string)
}

// DoTimestamp records the timestamp to the context with the given description.
// If the context does not implement TimestampContext, nothing happens.
// TODO 入れ替え
func DoTimestamp(ctx context.Context, description string) {
	if ctx, ok := ctx.(TimestampContext); ok {
		ctx.Timestamp(description)
	}
}

type timestampContext struct {
	context.Context
	createdAt time.Time
	maxStamps int

	mu     sync.Mutex
	stamps []elapsed
}

type elapsed struct {
	description string
	elapsed     time.Duration
}

type timestampContextOptions struct {
	maxStamps int
	since     time.Time
}

// TimestampContextOptionsFunc is a type of functional options for the TimeStampContext
type TimestampContextOptionsFunc func(*timestampContextOptions)

// WithTimestamp creates the TimestampContext
func WithTimestamp(parent context.Context, opts ...TimestampContextOptionsFunc) {
	conf := timestampContextOptions{
		maxStamps: 30,
		since:     time.Now(),
	}
	for _, o := range opts {
		o(&conf)
	}

}

func (c *timestampContext) Timestamp(description string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stamps = append(c.stamps, elapsed{
		description: description,
		elapsed:     time.Now().Sub(c.createdAt),
	})
}

package context

import (
	"context"
	"sync"
	"time"
)

// TODO test

// TimestampContext is the context in which time stamps can be recorded
type TimestampContext interface {
	context.Context
	// DoTimestamp records the timestamp to this context with the given description
	DoTimestamp(description string)
	// ListTimestamps lists timestamps in this context
	ListTimestamps() []Timestamp
}

// Timestamp represents a timestamp
type Timestamp struct {
	Description string
	Time        time.Time
}

// DoTimestamp records the timestamp to the context with the given description.
// If the context does not implement TimestampContext, nothing happens.
func DoTimestamp(ctx context.Context, description string) {
	if ctx, ok := ctx.(TimestampContext); ok {
		ctx.DoTimestamp(description)
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

// TimestampMax configures timestampContext.maxStamps
func TimestampMax(max int) TimestampContextOptionsFunc {
	return func(o *timestampContextOptions) {
		o.maxStamps = max
	}
}

// TimestampSince configures timestampContext.since
func TimestampSince(since time.Time) TimestampContextOptionsFunc {
	return func(o *timestampContextOptions) {
		o.since = since
	}
}

// WithTimestamp creates the TimestampContext
func WithTimestamp(parent context.Context, opts ...TimestampContextOptionsFunc) TimestampContext {
	conf := timestampContextOptions{
		maxStamps: 30,
		since:     time.Now(),
	}
	for _, o := range opts {
		o(&conf)
	}
	c := timestampContext{
		Context:   parent,
		createdAt: conf.since,
		maxStamps: conf.maxStamps,
		stamps:    make([]elapsed, 0, conf.maxStamps),
	}
	return &c
}

func (c *timestampContext) DoTimestamp(description string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.stamps) < c.maxStamps {
		c.stamps = append(c.stamps, elapsed{
			description: description,
			elapsed:     time.Now().Sub(c.createdAt),
		})
	}
}

func (c *timestampContext) ListTimestamps() []Timestamp {
	c.mu.Lock()
	defer c.mu.Unlock()

	ret := make([]Timestamp, len(c.stamps))
	for i, s := range c.stamps {
		ret[i] = Timestamp{
			Description: s.description,
			Time:        c.createdAt.Add(s.elapsed),
		}
	}
	return ret
}

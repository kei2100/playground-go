package time

import (
	"sync"
	"time"
)

// ChangeableTicker is a ticker that can change the interval
type ChangeableTicker struct {
	mu     sync.RWMutex
	ticker *time.Ticker
}

// NewChangeableTicker creates a ChangeableTicker
func NewChangeableTicker(d time.Duration) *ChangeableTicker {
	return &ChangeableTicker{
		ticker: time.NewTicker(d),
	}
}

// Change the interval
func (t *ChangeableTicker) Change(d time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ticker.Stop()
	t.ticker = time.NewTicker(d)
}

// C returns a channel on which the ticks are delivered.
func (t *ChangeableTicker) C() <-chan time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.ticker.C
}

// Stop turns off a ticker. After Stop, no more ticks will be sent.
// Stop does not close the channel, to prevent a read from the channel succeeding
// incorrectly.
func (t *ChangeableTicker) Stop() {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.ticker.Stop()
}

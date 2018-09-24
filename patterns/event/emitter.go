package event

import (
	"sync"
)

// Emitter is an event emitter
type Emitter struct {
	mu             sync.RWMutex
	topicListeners map[string]listeners
}

type listeners []chan interface{}

// Emit emits an event
func (e *Emitter) Emit(topic string, value interface{}) (done chan struct{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	done = make(chan struct{})

	if e.topicListeners == nil {
		close(done)
		return done
	}
	lns, ok := e.topicListeners[topic]
	if !ok || len(lns) == 0 {
		close(done)
		return done
	}

	go func() {
		defer close(done)
		for _, lnch := range lns {
			lnch <- value
		}
	}()
	return done
}

// On returns a channel that receives events
func (e *Emitter) On(topic string) <-chan interface{} {
	ch := make(chan interface{})

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.topicListeners == nil {
		e.topicListeners = make(map[string]listeners, 1)
		e.topicListeners[topic] = listeners{ch}
		return ch
	}

	if lns, ok := e.topicListeners[topic]; ok {
		e.topicListeners[topic] = append(lns, ch)
	} else {
		e.topicListeners[topic] = listeners{ch}
	}
	return ch
}

// Off removes ch from event listeners and closes ch.
// If the ch is not a listener, nothing happens.
func (e *Emitter) Off(topic string, ch <-chan interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.topicListeners == nil {
		return
	}
	lns, ok := e.topicListeners[topic]
	if !ok {
		return
	}
	for i, lnch := range lns {
		if lnch == ch {
			e.topicListeners[topic] = append(lns[:i], lns[i+1:]...)
			close(lnch)
		}
	}
}

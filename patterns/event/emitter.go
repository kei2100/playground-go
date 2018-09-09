package event

import (
	"reflect"
	"sync"
)

// AEvent hoghoge is an event
type AEvent struct{}

type emitter struct {
	mu sync.RWMutex

	listenersA map[uintptr]struct {
		cb   func(event *AEvent)
		remv func()
	}
}

// EmitA event
func (e *emitter) EmitA(ev *AEvent) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, ln := range e.listenersA {
		go ln.cb(ev)
	}
}

// OnA set the callback to be called when an event emitted
func (e *emitter) OnA(callback func(event *AEvent)) (remove func()) {
	ptr := reflect.ValueOf(callback).Pointer()

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.listenersA == nil {
		e.listenersA = make(map[uintptr]struct {
			cb   func(event *AEvent)
			remv func()
		})
	}
	if ln, ok := e.listenersA[ptr]; ok {
		return ln.remv
	}

	var once sync.Once
	remv := func() {
		once.Do(func() {
			e.mu.Lock()
			defer e.mu.Unlock()
			delete(e.listenersA, ptr)
		})
	}
	e.listenersA[ptr] = struct {
		cb   func(event *AEvent)
		remv func()
	}{cb: callback, remv: remv}

	return remv
}

package event

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

type eventType reflect.Type

// Handler is a type of event handler function.
// Must have one argument representing a type of the event
type Handler interface{}
type handlers map[uintptr]Handler

// UntypedEmitter is a reflection based event emitter
type UntypedEmitter struct {
	mu          sync.RWMutex
	allHandlers map[eventType]handlers
}

// On adds handlerFunc to the event handlers.
// The handlerFunc must have one argument representing a type of the event.
func (e *UntypedEmitter) On(handlerFunc Handler) {
	hv, et, err := reflectHandler(handlerFunc)
	if err != nil {
		log.Println(err)
		return
	}
	hp := hv.Pointer()

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.allHandlers == nil {
		e.allHandlers = make(map[eventType]handlers)
		e.allHandlers[et] = handlers{hp: handlerFunc}
		return
	}
	if e.allHandlers[et] == nil {
		e.allHandlers[et] = handlers{hp: handlerFunc}
		return
	}
	e.allHandlers[et][hp] = handlerFunc
}

// Off removes handlerFunc from event handlers.
// If the handlerFunc not registered, it has no effect
func (e *UntypedEmitter) Off(handlerFunc Handler) {
	hv, et, err := reflectHandler(handlerFunc)
	if err != nil {
		log.Println(err)
		return
	}
	hp := hv.Pointer()

	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.allHandlers) == 0 {
		return
	}
	if len(e.allHandlers[et]) == 0 {
		return
	}
	delete(e.allHandlers[et], hp)
	if len(e.allHandlers[et]) == 0 {
		delete(e.allHandlers, et)
	}
}

// Handlers returns handlers
func (e *UntypedEmitter) Handlers(et eventType) []Handler {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if len(e.allHandlers) == 0 {
		return nil
	}
	hns := e.allHandlers[et]
	if len(hns) == 0 {
		return nil
	}
	ret := make([]Handler, len(hns))
	i := 0
	for _, h := range hns {
		ret[i] = h
		i++
	}
	return ret
}

// Emit the event
func (e *UntypedEmitter) Emit(ev interface{}) {
	evv := reflect.ValueOf(ev)
	et := evv.Type()
	hns := e.Handlers(et)

	for _, h := range hns {
		handlerFunc := h
		hv, argt, err := reflectHandler(handlerFunc)
		if err != nil {
			log.Println(err)
			continue
		}
		if et != argt {
			log.Printf("event: event type and arg type is not same %T, %T", ev, handlerFunc)
		}
		hv.Call([]reflect.Value{evv})
	}
}

func reflectHandler(handlerFunc Handler) (handlerValue reflect.Value, eventType eventType, err error) {
	hv := reflect.ValueOf(handlerFunc)
	if hv.Kind() != reflect.Func {
		log.Println("event: handler is not a function")
		return reflect.Value{}, nil, fmt.Errorf("event: handler %T is not a function", handlerFunc)
	}
	ht := hv.Type()
	if ht.NumIn() != 1 {
		return reflect.Value{}, nil, fmt.Errorf("event: handler %T must have one argument", handlerFunc)
	}
	et := ht.In(0)
	return hv, et, nil
}

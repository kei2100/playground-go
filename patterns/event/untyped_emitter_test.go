package event

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

type aEvent struct {
	msg string
}
type bEvent struct {
	msg string
}

func TestUntypedEmitter_On(t *testing.T) {
	t.Run("handler kinds", func(t *testing.T) {
		tt := []struct {
			handlerFunc  interface{}
			wantRegister bool
		}{
			{handlerFunc: nil, wantRegister: false},
			{handlerFunc: "", wantRegister: false},
			{handlerFunc: func() {}, wantRegister: false},
			{handlerFunc: func(ev *aEvent) {}, wantRegister: true},
			{handlerFunc: func(e1 *aEvent, e2 *bEvent) {}, wantRegister: false},
		}
		for _, te := range tt {
			e := new(UntypedEmitter)
			e.On(te.handlerFunc)
			if te.wantRegister {
				if g, w := len(e.allHandlers), 0; !(g > w) {
					t.Errorf("allHandlers length got %v, want > %v", g, w)
				}
			} else {
				if g, w := len(e.allHandlers), 0; g != w {
					t.Errorf("allHandlers length got %v, want %v", g, w)
				}
			}
		}
	})

	t.Run("register some handlers", func(t *testing.T) {
		e := new(UntypedEmitter)

		aHandler := func(_ *aEvent) {}
		bHandler := func(_ *bEvent) {}
		aHandler2 := func(_ *aEvent) {}

		wg := sync.WaitGroup{}
		wg.Add(4)
		go func() { e.On(aHandler); wg.Done() }()
		go func() { e.On(bHandler); wg.Done() }()
		go func() { e.On(aHandler2); wg.Done() }()
		go func() { e.On(aHandler); wg.Done() }()
		wg.Wait()

		if e.allHandlers == nil {
			t.Fatal("allHandlers is nil")
		}
		if g, w := len(e.allHandlers[reflect.TypeOf(&aEvent{})]), 2; g != w {
			t.Errorf("aEvent handlers lentgh got %v, want %v", g, w)
		}
		if g, w := len(e.allHandlers[reflect.TypeOf(&bEvent{})]), 1; g != w {
			t.Errorf("bEvent handlers lentgh got %v, want %v", g, w)
		}
	})
}

func TestUntypedEmitter_Off(t *testing.T) {
	e := new(UntypedEmitter)

	e.Off(nil)
	e.Off("")
	e.Off(func() {})
	e.Off(func(_ *aEvent, _ *bEvent) {})

	if len(e.allHandlers) != 0 {
		t.Fatal("allHandlers is not empty")
	}

	aHandler := func(_ *aEvent) {}
	bHandler := func(_ *bEvent) {}
	aHandler2 := func(_ *aEvent) {}

	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() { e.On(aHandler); wg.Done() }()
	go func() { e.On(bHandler); wg.Done() }()
	go func() { e.On(aHandler2); wg.Done() }()
	go func() { e.On(aHandler); wg.Done() }()
	wg.Wait()

	wg.Add(2)
	go func() { e.Off(aHandler); wg.Done() }()
	go func() { e.Off(bHandler); wg.Done() }()
	wg.Wait()

	if g, w := len(e.allHandlers), 1; g != w {
		t.Fatalf("allHandlers length got %v, want %v", g, w)
	}
	hns, ok := e.allHandlers[reflect.TypeOf(&aEvent{})]
	if !ok {
		t.Error("aEvent handler is not exist")
	}
	if _, ok := hns[reflect.ValueOf(aHandler2).Pointer()]; !ok {
		t.Error("aHandler2 is not exist")
	}
}

func TestUntypedEmitter_Emit(t *testing.T) {
	e := new(UntypedEmitter)

	var aResult, bResult, aResult2 string
	wg := sync.WaitGroup{}
	aHandler := func(ev *aEvent) { aResult = ev.msg; wg.Done() }
	bHandler := func(ev *bEvent) { bResult = ev.msg; wg.Done() }
	aHandler2 := func(ev *aEvent) { aResult2 = ev.msg; wg.Done() }

	e.On(aHandler)
	e.On(bHandler)
	e.On(aHandler2)

	wg.Add(2)
	e.Emit(&aEvent{msg: "a"})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		break
	case <-time.After(time.Second):
		t.Errorf("timeout exceeded while waiting for handle")
	}

	if g, w := aResult, "a"; g != w {
		t.Errorf("aResult got %v, want %v", g, w)
	}
	if g, w := aResult2, "a"; g != w {
		t.Errorf("aResult2 got %v, want %v", g, w)
	}
	if g, w := bResult, ""; g != w {
		t.Errorf("bResult got %v, want %v", g, w)
	}
}

const nHandlers = 10

// Benchmark_UseReflection-8        3000000               401 ns/op
func Benchmark_UseReflection(b *testing.B) {
	e := new(UntypedEmitter)
	for i := 0; i < nHandlers; i++ {
		e.On(func(_ *aEvent) {})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e.Emit(&aEvent{})
	}
}

// Benchmark_UseTypeAssertion-8    10000000               165 ns/op
func Benchmark_UseTypeAssertion(b *testing.B) {
	e := new(UntypedEmitter)
	for i := 0; i < nHandlers; i++ {
		e.On(func(_ *aEvent) {})
	}

	b.ResetTimer()
	ev := &aEvent{}

	for i := 0; i < b.N; i++ {
		hns := e.Handlers(reflect.TypeOf(ev))
		for _, h := range hns {
			if h, ok := h.(func(*aEvent)); ok {
				h(ev)
			}
		}
	}
}

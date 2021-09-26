package event

import (
	"reflect"
	"testing"
	"time"
)

func TestEmitter_On_Emit(t *testing.T) {
	em := Emitter{}

	foo1 := em.On("foo")
	foo2 := em.On("foo")
	bar := em.On("bar")

	var fooval1, fooval2 string

	done := em.Emit("foo", "hoo")
	testNotReceiveStruct(t, done)

	val1 := testReceive(t, foo1, time.Second)
	assignIfString(&fooval1, val1)
	testNotReceiveStruct(t, done)

	val2 := testReceive(t, foo2, time.Second)
	assignIfString(&fooval2, val2)
	testReceiveStruct(t, done, time.Second)

	testNotReceive(t, bar)

	if g, w := fooval1, "hoo"; g != w {
		t.Errorf("fooval1 got %v, want %v", g, w)
	}
	if g, w := fooval2, "hoo"; g != w {
		t.Errorf("fooval2 got %v, want %v", g, w)
	}
}

func TestEmitter_Off_Emit(t *testing.T) {
	em := Emitter{}

	foo1 := em.On("foo")
	foo2 := em.On("foo")
	foo3 := em.On("foo")

	em.Off("foo", foo2)
	var fooval1, fooval2, fooval3 string

	done := em.Emit("foo", "hoo")
	testNotReceiveStruct(t, done)

	val1 := testReceive(t, foo1, time.Second)
	assignIfString(&fooval1, val1)
	testNotReceiveStruct(t, done)

	val2 := testReceive(t, foo2, time.Second)
	assignIfString(&fooval2, val2)
	testNotReceiveStruct(t, done)

	val3 := testReceive(t, foo3, time.Second)
	assignIfString(&fooval3, val3)
	testReceiveStruct(t, done, time.Second)

	if g, w := fooval1, "hoo"; g != w {
		t.Errorf("fooval1 got %v, want %v", g, w)
	}
	if g, w := fooval2, ""; g != w {
		t.Errorf("fooval2 got %v, want %v", g, w)
	}
	if g, w := fooval3, "hoo"; g != w {
		t.Errorf("fooval3 got %v, want %v", g, w)
	}
}

func TestEmitter_Emit_NoListener(t *testing.T) {
	em := Emitter{}

	d1 := em.Emit("foo", "foo")
	testReceiveStruct(t, d1, time.Second)

	em.On("bar")

	d2 := em.Emit("foo", "foo")
	testReceiveStruct(t, d2, time.Second)
}

func testNotReceiveStruct(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
		t.Error("unexpected ch receive")
	default:
	}
}

func testReceiveStruct(t *testing.T, ch <-chan struct{}, timeout time.Duration) {
	t.Helper()
	select {
	case <-ch:
		return
	case <-time.After(timeout):
		t.Errorf("timeout %s while waiting for ch receive", timeout)
		return
	}
}

func testNotReceive(t *testing.T, ch <-chan interface{}) {
	t.Helper()
	select {
	case <-ch:
		t.Error("unexpected ch receive")
	default:
	}
}

func testReceive(t *testing.T, ch <-chan interface{}, timeout time.Duration) interface{} {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(timeout):
		t.Errorf("timeout %s while waiting for ch receive", timeout)
		return nil
	}
}

func assignIfString(to *string, from interface{}) {
	rv := reflect.ValueOf(from)
	if rv.Kind() != reflect.String {
		return
	}
	*to = rv.String()
}

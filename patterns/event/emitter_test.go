package event

import (
	"testing"
	"time"
)

var ac chan struct{}

func privFuncA(_ *AEvent) {
	close(ac)
}

func GlobFuncA(_ *AEvent) {
	close(ac)
}

func TestEmitter(t *testing.T) {
	ac = make(chan struct{})

	localFuncA := func(_ *AEvent) {
		close(ac)
	}

	e := new(emitter)
	e.OnA(privFuncA)
	e.OnA(GlobFuncA)
	e.OnA(localFuncA)
	// check for ignore duplicate registration
	removeAPriv := e.OnA(privFuncA)
	removeAGlob := e.OnA(GlobFuncA)
	removeALocal := e.OnA(localFuncA)

	removeAPriv()
	removeAGlob()
	defer removeALocal()

	e.EmitA(&AEvent{})

	select {
	case <-ac:
		break
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout exceeded while waiting for the channel")
	}
}

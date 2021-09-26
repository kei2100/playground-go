package reflect

import (
	"reflect"
	"testing"
)

func TestReflectChannel(t *testing.T) {
	ch1 := make(chan<- string, 1)
	rv1 := reflect.ValueOf(ch1)

	if rv1.Kind() != reflect.Chan {
		t.Fatalf("unexpected type %s", rv1.Kind())
	}

	if rv1.Cap() != 1 {
		t.Fatalf("unexpected cap %d", rv1.Cap())
	}
	ch1 <- ""
	if rv1.Cap() != 1 { // 変わらない
		t.Fatalf("unexpected cap %d", rv1.Cap())
	}

	rt1 := rv1.Type()
	if rt1.ChanDir() != reflect.SendDir {
		t.Fatalf("unexpected direction %s", rt1.ChanDir())
	}

	if rt1.Elem().Kind() != reflect.String {
		t.Fatalf("unexpected elem type %s", rt1.Elem().Kind())
	}
}

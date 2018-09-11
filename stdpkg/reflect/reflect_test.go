package reflect

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestCheckImplementsInterface(t *testing.T) {
	err := errors.New("test")
	et := reflect.TypeOf(err)

	ok := et.Implements(reflect.TypeOf((*error)(nil)).Elem())
	if !ok {
		t.Error("not ok")
	}

	ok = et.Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem())
	if ok {
		t.Error("want not ok")
	}
}

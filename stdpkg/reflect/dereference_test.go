package reflect

import (
	"reflect"
	"testing"
)

func TestDereference(t *testing.T) {
	type s struct {
		Message string
	}
	sv := s{Message: "test"}

	it := reflect.ValueOf(sv).Interface()
	it1 := dereference(sv).Interface()
	it2 := dereference(&sv).Interface()

	if it != it1 {
		t.Errorf("not same %T, %T", it, it1)
	}
	if it != it2 {
		t.Errorf("not same %T, %T", it, it2)
	}
}

func dereference(v interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(v))
}

package reflect

import (
	"log"
	"reflect"
	"testing"
)

func TestReference(t *testing.T) {
	// T => *T への変換

	type myStruct struct {
	}
	var ms interface{} = myStruct{}

	v := reflect.ValueOf(ms)
	pt := reflect.PtrTo(v.Type()) // create a *T type.
	pv := reflect.New(pt.Elem())  // create a reflect.Value of type *T.
	pv.Elem().Set(v)              // sets pv to point to underlying value of v.

	log.Printf("%T", pv.Interface()) // *reflect.myStruct
}

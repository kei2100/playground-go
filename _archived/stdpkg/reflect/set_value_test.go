package reflect

import (
	"log"
	"reflect"
	"testing"
)

func TestSetToInterface(t *testing.T) {
	type MyInterface interface{}
	type MyImpl struct {
		MyImplString string
	}

	type MyStruct struct {
		MyString string
		MyType1  MyInterface
		MyType2  MyInterface
	}

	src := &MyStruct{
		MyString: "string value",
		MyType1: &MyImpl{
			MyImplString: "myimpl string value",
		},
	}
	var dest MyStruct

	// src *MyStruct
	srv := reflect.ValueOf(src)
	// src MyStruct
	srv = reflect.Indirect(srv)

	// dest *MyStruct
	drv := reflect.ValueOf(&dest)

	for i := 0; i < srv.NumField(); i++ {
		sfv := srv.Field(i)
		if sfv.Kind() != reflect.Interface {
			continue
		}
		if sfv.IsNil() {
			continue
		}
		sft := reflect.PtrTo(sfv.Type())
		// dfv *MyImpl
		dfv := reflect.New(sft.Elem())
		// dfvが参照している値にsfvをセット
		dfv.Elem().Set(sfv)

		// dest *MyStruct が参照している値のフィールドにdfvの値をセット
		drv.Elem().Field(i).Set(dfv.Elem())
	}

	log.Printf("%+v", dest)
}

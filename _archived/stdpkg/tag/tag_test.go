package tag

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type MyStruct struct {
	Foo string `tagA:"foo-a" tagB:"foo-b"`
}

func TestReflectStructTags(t *testing.T) {
	s := MyStruct{}
	sp := &MyStruct{}
	var snp *MyStruct

	if err := reflectStructTags(s); err != nil {
		t.Error(err)
	}
	if err := reflectStructTags(sp); err != nil {
		t.Error(err)
	}
	if err := reflectStructTags(snp); err != nil {
		t.Error(err)
	}
}

func reflectStructTags(st interface{}) error {
	rt := reflect.TypeOf(st)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return errors.New("not struct")
	}
	nf := rt.NumField()
	for i := 0; i < nf; i++ {
		f := rt.Field(i)
		ta := f.Tag.Get("tagA")
		fmt.Println(ta)
		tb := f.Tag.Get("tagB")
		fmt.Println(tb)
	}
	return nil
}

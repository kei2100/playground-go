package reflect

import (
	"log"
	"reflect"
	"testing"
)

func TestValueOf(t *testing.T) {
	{
		type myStruct struct {
		}
		var v *myStruct
		rv := reflect.ValueOf(v)
		log.Printf("kind:%v isnil:%v", rv.Kind(), rv.IsNil())
	}
	{
		v := interface{}(interface{}("a"))
		rv := reflect.ValueOf(v)
		log.Println(rv.Kind()) // string
	}
	{
		type myString string
		v := myString("my")
		rv := reflect.ValueOf(v)
		log.Println(rv.Kind()) // string
	}
}

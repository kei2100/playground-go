package reflect

import (
	"log"
	"reflect"
	"testing"
)

func TestKind(t *testing.T) {
	{
		type myStruct struct{}
		var v *myStruct
		rv := reflect.ValueOf(v)
		log.Printf("kind:%v isnil:%v", rv.Kind(), rv.IsNil()) // kind:ptr isnil:true
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
	{
		v := make(map[string]interface{})
		rv := reflect.ValueOf(v)
		log.Println(rv.Kind()) // map
	}
	{
		type mymap map[string]interface{}
		var v mymap
		rv := reflect.ValueOf(v)
		log.Println(rv.Kind()) // map
	}
	{
		v := map[string]interface{}{}
		rv := reflect.ValueOf(v)
		log.Println(rv.Kind()) // map
	}
	{
		v := map[string]interface{}{
			"vv": map[string]interface{}{
				"vvs": "string",
			},
		}
		rv := reflect.ValueOf(v["vv"])
		log.Println(rv.Kind()) // map
	}
	{
		type myStruct struct{}
		var m map[string]*myStruct
		rt := reflect.TypeOf(m).Elem()
		log.Println(rt) // *reflect.myStruct
	}
	{
		type myStruct struct{}
		var m map[string]*myStruct
		rt := reflect.TypeOf(m).Key()
		log.Println(rt) // string
	}
}

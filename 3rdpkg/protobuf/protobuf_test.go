package protobuf

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"unicode"

	structpb "github.com/golang/protobuf/ptypes/struct"
)

func EncodeFromMap(m map[string]interface{}) (*structpb.Struct, error) {
	var pbst structpb.Struct
	pbst.Fields = make(map[string]*structpb.Value, len(m))
	for k, v := range m {
		pbv, err := EncodeValue(v)
		if err != nil {
			return nil, err
		}
		pbst.Fields[k] = pbv
	}
	return &pbst, nil
}

func EncodeFromStruct(v interface{}) (*structpb.Struct, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, nil
		}
		rv = reflect.Indirect(rv)
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("oops")
	}
	rt := rv.Type()
	pbst := structpb.Struct{Fields: make(map[string]*structpb.Value)}
	for i := 0; i < rt.NumField(); i++ {
		name := rt.Field(i).Name
		if !unicode.IsUpper(rune(name[0])) {
			continue
		}
		pbvalue, err := EncodeValue(rv.FieldByName(name).Interface())
		if err != nil {
			return nil, err
		}
		pbst.Fields[name] = pbvalue
	}
	return &pbst, nil
}

func EncodeValue(v interface{}) (*structpb.Value, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return &structpb.Value{Kind: &structpb.Value_NullValue{}}, nil
		}
		rv = reflect.Indirect(rv)
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(rv.Int())}}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(rv.Uint())}}, nil
	case reflect.Float32, reflect.Float64:
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: rv.Float()}}, nil
	case reflect.String:
		return &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: rv.String()}}, nil
	case reflect.Bool:
		return &structpb.Value{Kind: &structpb.Value_BoolValue{BoolValue: rv.Bool()}}, nil
	case reflect.Struct:
		pbst, err := EncodeFromStruct(rv.Interface())
		if err != nil {
			return nil, err
		}
		return &structpb.Value{Kind: &structpb.Value_StructValue{StructValue: pbst}}, nil

	case reflect.Map:
		keys := rv.MapKeys()
		pbst := structpb.Struct{Fields: make(map[string]*structpb.Value, len(keys))}
		for _, kv := range keys {
			if kv.Kind() != reflect.String {
				return nil, fmt.Errorf("unknown key type: %v", rv.Kind())
			}
			pbv, err := EncodeValue(rv.MapIndex(kv).Interface())
			if err != nil {
				return nil, err
			}
			pbst.Fields[kv.String()] = pbv
		}
		return &structpb.Value{Kind: &structpb.Value_StructValue{StructValue: &pbst}}, nil

	case reflect.Slice, reflect.Array:
		rvlen := rv.Len()
		pblv := structpb.ListValue{Values: make([]*structpb.Value, rvlen)}
		for i := 0; i < rvlen; i++ {
			pbv, err := EncodeValue(rv.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			pblv.Values[i] = pbv
		}
		return &structpb.Value{Kind: &structpb.Value_ListValue{ListValue: &pblv}}, nil
	}

	return nil, fmt.Errorf("unknown: %v", rv.Kind())
}

func TestEncode(t *testing.T) {
	//sp := "ss"
	//var snp *string
	m := map[string]interface{}{
		//"1":      0.1,
		//"2":      1,
		//"str":    "str",
		//"strPtr": &sp,
		//"nil":    snp,
		//"bool": true,
		//"map": make(map[string]interface{}),
		"map": map[string]interface{}{
			"foo": "bar",
		},
		"struct": struct {
			Foo string
			bar string
		}{
			Foo: "foo",
			bar: "bar",
		},
		"arr":   [2]string{"foo", "bar"},
		"slice": []string{"foo", "bar"},
	}
	s, err := EncodeFromMap(m)
	log.Printf("\ns:%+v\nerr:%v", s, err)
}

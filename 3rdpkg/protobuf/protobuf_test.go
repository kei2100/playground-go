package protobuf

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/struct"
	"log"
	"reflect"
	"testing"
)

func EncodeFromMap(m map[string]interface{}) (*structpb.Struct, error) {
	var st structpb.Struct
	st.Fields = make(map[string]*structpb.Value, len(m))
	for k, v := range m {
		pv, err := EncodeValue(v)
		if err != nil {
			return nil, err
		}
		st.Fields[k] = pv
	}
	return &st, nil
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
		// TODO
	case reflect.Map:
		keys := rv.MapKeys()
		pbstruct := structpb.Struct{Fields: make(map[string]*structpb.Value, len(keys))}
		for _, kv := range rv.MapKeys() {
			if kv.Kind() != reflect.String {
				return nil, fmt.Errorf("unknown key type: %v", rv.Kind())
			}
			pbvalue, err := EncodeValue(rv.MapIndex(kv))
			if err != nil {
				return nil, err
			}
			pbstruct.Fields[kv.String()] = pbvalue
		}
		return &structpb.Value{Kind: &structpb.Value_StructValue{StructValue: &pbstruct}}, nil
	case reflect.Slice, reflect.Array:
		// TODO
	}
	return nil, fmt.Errorf("unknown: %v", rv.Kind())
}

func TestEncode(t *testing.T) {
	sp := "ss"
	var snp *string
	m := map[string]interface{}{
		"1":      0.1,
		"2":      1,
		"str":    "str",
		"strPtr": &sp,
		"nil":    snp,
	}
	s, err := EncodeFromMap(m)
	log.Printf("\ns:%+v\nerr:%v", s, err)
}

package gateway

import (
	"reflect"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

// FilterFunc type
type FilterFunc func(tagValue string) bool

// TagBasedRequestFieldFilter -
func TagBasedRequestFieldFilter(tagName string, fn FilterFunc) grpc_ctxtags.RequestFieldExtractorFunc {
	return func(fullMethod string, req interface{}) map[string]interface{} {
		retMap := make(map[string]interface{})
		reflectMessageTags(req, retMap, tagName, fn)
		if len(retMap) == 0 {
			return nil
		}
		return retMap
	}
}

func reflectMessageTags(msg interface{}, existingMap map[string]interface{}, tagName string, fn FilterFunc) {
	v := reflect.ValueOf(msg)
	// Only deal with pointers to structs.
	if v.Kind() != reflect.Ptr {
		return
	}
	// Deref the pointer get to the struct.
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		kind := field.Kind()
		// Only recurse down direct pointers, which should only be to nested structs.
		if kind == reflect.Ptr {
			reflectMessageTags(field.Interface(), existingMap, tagName, fn)
		}
		// In case of arrays/slices (repeated fields) go down to the concrete type.
		if kind == reflect.Array || kind == reflect.Slice {
			if field.Len() == 0 {
				continue
			}
			kind = field.Index(0).Kind()
		}
		// Only be interested in
		if (kind >= reflect.Bool && kind <= reflect.Float64) || kind == reflect.String {
			t := v.Type()
			fieldType := t.Field(i)
			tagValue := fieldType.Tag.Get(tagName)
			fieldName := fieldType.Name
			if fn(tagValue) {
				existingMap[fieldName] = field.Interface()
			}
		}
	}
	return
}

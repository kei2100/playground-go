package reflect

import (
	"fmt"
	"reflect"
	"testing"
)

func TestIsNilPointer(t *testing.T) {
	var isNil = func(v interface{}) bool {
		rv := reflect.ValueOf(v)
		return rv.Kind() == reflect.Ptr && rv.Pointer() == 0
	}

	var strp *string

	tt := []struct {
		v    interface{}
		want bool
	}{
		{
			v:    nil,
			want: false,
		},
		{
			v:    strp,
			want: true,
		},
	}
	for i, te := range tt {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			if got := isNil(te.v); got != te.want {
				t.Errorf("got %v, want %v", got, te.want)
			}
		})
	}
}

package gateway

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagBasedRequestFieldFilter(t *testing.T) {
	tt := []struct {
		in   interface{}
		want map[string]interface{}
	}{
		{
			in: &struct {
				A string `log:"-"`
				B string
			}{
				A: "a",
				B: "b",
			},
			want: map[string]interface{}{
				"B": "b",
			},
		},
	}
	for i, te := range tt {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			f := TagBasedRequestFieldFilter("log", func(tagValue string) bool { return tagValue != "-" })
			got := f("", te.in)
			assert.Equal(t, te.want, got)
		})
	}
}

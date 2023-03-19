package slice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteAt(t *testing.T) {
	tt := []struct {
		title string
		s     []int
		i     int
		want  []int
	}{
		{
			title: "works",
			s:     []int{0, 1, 2, 3, 4},
			i:     2,
			want:  []int{0, 1, 3, 4},
		},
		{
			title: "first",
			s:     []int{0, 1, 2, 3, 4},
			i:     0,
			want:  []int{1, 2, 3, 4},
		},
		{
			title: "last",
			s:     []int{0, 1, 2, 3, 4},
			i:     4,
			want:  []int{0, 1, 2, 3},
		},
	}
	for i, te := range tt {
		t.Run(fmt.Sprintf("#%d %s", i, te.title), func(t *testing.T) {
			got := DeleteAt(te.s, te.i)
			assert.Equal(t, te.want, got)
		})
	}
}

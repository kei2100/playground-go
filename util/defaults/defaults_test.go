package defaults

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	zero := time.Time{}
	v1 := time.Unix(1, 0)
	v2 := time.Unix(2, 0)

	tests := []struct {
		v        time.Time
		defaultV time.Time
		wantV    time.Time
	}{
		{zero, zero, zero},
		{zero, v1, v1},
		{v2, v1, v2},
	}

	for i, tt := range tests {
		if g, w := Time(tt.v, tt.defaultV), tt.wantV; g != w {
			t.Errorf("tests[%v] got %v, want %v", i, g, w)
		}
	}
}

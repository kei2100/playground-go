package slice

import (
	"testing"
)

// $ go test -bench BenchmarkAppendAlloc -benchmem bench/slice/slice_test.go
// BenchmarkAppendAlloc/cap_0_append_100-4                  2000000               779 ns/op            2040 B/op          8 allocs/op
// BenchmarkAppendAlloc/cap_100_append_100-4                5000000               298 ns/op             896 B/op          1 allocs/op
// BenchmarkAppendAlloc/cap_0_append_10000-4                  20000             78105 ns/op          386296 B/op         20 allocs/op
// BenchmarkAppendAlloc/cap_10000_append_10000-4              50000             28455 ns/op           81920 B/op          1 allocs/op
func BenchmarkAppendAlloc(b *testing.B) {
	tests := []struct {
		subject string
		cap     int
		appendN int
	}{
		{
			subject: "cap 0 append 100",
			cap:     0,
			appendN: 100,
		},
		{
			subject: "cap 100 append 100",
			cap:     100,
			appendN: 100,
		},
		{
			subject: "cap 0 append 10000",
			cap:     0,
			appendN: 10000,
		},
		{
			subject: "cap 10000 append 10000",
			cap:     10000,
			appendN: 10000,
		},
	}
	for _, tt := range tests {
		b.Run(tt.subject, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s := make([]int, 0, tt.cap)
				for j := 0; j < tt.appendN; j++ {
					s = append(s, j)
				}
			}
		})
	}
}

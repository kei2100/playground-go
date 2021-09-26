package slice

import (
	"testing"
)

// $ go test -bench BenchmarkAppendAllocCap -benchmem bench/slice/slice_test.go
// BenchmarkAppendAllocCap/cap_0_append_100-4                  2000000               779 ns/op            2040 B/op          8 allocs/op
// BenchmarkAppendAllocCap/cap_100_append_100-4                5000000               298 ns/op             896 B/op          1 allocs/op
// BenchmarkAppendAllocCap/cap_0_append_10000-4                  20000             78105 ns/op          386296 B/op         20 allocs/op
// BenchmarkAppendAllocCap/cap_10000_append_10000-4              50000             28455 ns/op           81920 B/op          1 allocs/op
func BenchmarkAppendAllocCap(b *testing.B) {
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

// $ go test -bench BenchmarkAppendAllocJoin -benchmem bench/slice/slice_test.go
// BenchmarkAppendAllocJoin/define_capped_destination_slice-4              100000000               11.7 ns/op             0 B/op          0 allocs/op
// BenchmarkAppendAllocJoin/define_uncapped_destination_slice-4            10000000               144 ns/op             224 B/op          3 allocs/op
// BenchmarkAppendAllocJoin/not_define_destination_slice-4                 20000000                95.4 ns/op           144 B/op          2 allocs/op
func BenchmarkAppendAllocJoin(b *testing.B) {
	b.Run("define capped destination slice", func(b *testing.B) {
		s1 := []int{0, 1, 2}
		s2 := []int{3, 4, 5}
		s3 := []int{6, 7, 8}

		for i := 0; i < b.N; i++ {
			s := make([]int, 0, 9)
			s = append(s, s1...)
			s = append(s, s2...)
			s = append(s, s3...)
		}
	})

	b.Run("define uncapped destination slice", func(b *testing.B) {
		s1 := []int{0, 1, 2}
		s2 := []int{3, 4, 5}
		s3 := []int{6, 7, 8}

		for i := 0; i < b.N; i++ {
			s := make([]int, 0)
			s = append(s, s1...) // alloc
			s = append(s, s2...) // alloc
			s = append(s, s3...) // alloc
		}
	})

	b.Run("not define destination slice", func(b *testing.B) {
		s1 := []int{0, 1, 2}
		s2 := []int{3, 4, 5}
		s3 := []int{6, 7, 8}

		for i := 0; i < b.N; i++ {
			s := append(s1, s2...) // alloc
			s = append(s, s3...)   // alloc
		}
	})
}

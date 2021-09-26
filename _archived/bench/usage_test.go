package bench

import "testing"

// $ go test -bench BenchmarkDefinePosition -benchmem bench/usage_test.go
// BenchmarkDefinePosition/Inside-4                50000000                29.3 ns/op             8 B/op          1 allocs/op
// BenchmarkDefinePosition/Outside-4               50000000                32.0 ns/op             8 B/op          1 allocs/op
func BenchmarkDefinePosition(b *testing.B) {
	b.Run("Inside", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := make([]int, 0)
			s = append(s, i) // alloc
		}
	})

	b.Run("Outside", func(b *testing.B) {
		s := make([]int, 0)
		s = append(s, 0) // alloc （だが、カウントされない。この部分はレコードされない模様）
		for i := 0; i < b.N; i++ {
			s = make([]int, 0)
			s = append(s, i) // alloc
		}
	})
}

package channel

import "testing"

// $ go test -bench BenchmarkChanHolder -benchmem
//  goos: darwin
//  goarch: amd64
//  pkg: github.com/kei2100/playground-go/src/benchmark/channel
//  cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
//  BenchmarkChanHolder1_OnceCloseChannel-16        257844824                4.678 ns/op           0 B/op          0 allocs/op
//  BenchmarkChanHolder2_OnceCloseChannel-16        781616865                1.498 ns/op           0 B/op          0 allocs/op

func BenchmarkChanHolder1_OnceCloseChannel(b *testing.B) {
	h := NewChanHolder1()
	for i := 0; i < b.N; i++ {
		h.OnceCloseChannel()
	}
}

func BenchmarkChanHolder2_OnceCloseChannel(b *testing.B) {
	h := NewChanHolder2()
	for i := 0; i < b.N; i++ {
		h.OnceCloseChannel()
	}
}

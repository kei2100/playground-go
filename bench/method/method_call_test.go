package method

import (
	"testing"
)

type valueReceiver struct {
	n int
}

func (r valueReceiver) incr(v int) int {
	return v + r.n
}

type pointerReceiver struct {
	n int
}

func (r *pointerReceiver) incr(v int) int {
	return v + r.n
}

// $ go test -bench BenchmarkMethodCall -benchmem bench/method/method_call_test.go
// BenchmarkMethodCall/value_receiver-4            2000000000               0.33 ns/op            0 B/op          0 allocs/op
// BenchmarkMethodCall/pointer_receiver-4          2000000000               0.43 ns/op            0 B/op          0 allocs/op
// BenchmarkMethodCall/implicit_pointer_receiver-4 2000000000               0.42 ns/op            0 B/op          0 allocs/op
func BenchmarkMethodCall(b *testing.B) {
	vr := valueReceiver{n: 1}
	pr := &pointerReceiver{n: 1}
	ipr := pointerReceiver{n: 1}


	b.Run("value receiver", func(b *testing.B) {
		i, v := 0, 0
		for ; i < b.N; i++ {
			v = vr.incr(v)
		}
		if g, w := v, i; g != w {
			b.Errorf("v got %v, want %v", g, w)
		}
	})
	b.Run("pointer receiver", func(b *testing.B) {
		i, v := 0, 0
		for ; i < b.N; i++ {
			v = pr.incr(v)
		}
		if g, w := v, i; g != w {
			b.Errorf("v got %v, want %v", g, w)
		}
	})
	b.Run("implicit pointer receiver", func(b *testing.B) {
		i, v := 0, 0
		for ; i < b.N; i++ {
			v = ipr.incr(v)
		}
		if g, w := v, i; g != w {
			b.Errorf("v got %v, want %v", g, w)
		}
	})
}

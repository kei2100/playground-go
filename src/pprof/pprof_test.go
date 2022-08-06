package pprof

import (
	"math"
	"reflect"
	"testing"

	"github.com/pkg/profile"
)

func TestPProf(t *testing.T) {
	const n = 100000
	// default: cpu profiling
	prof := profile.Start(
		profile.ProfilePath("./testdata/"),
		profile.NoShutdownHook,
	)
	result1 := primalityTestSlow(n)
	result2 := primalityTestFast(n)
	prof.Stop()

	// view profile
	// $ go tool pprof -http=":8081"  ./testdata/cpu.pprof

	if !reflect.DeepEqual(result1, result2) {
		t.Errorf("not same result1, result2\n%v\n%v\n", result1, result2)
	}
}

func primalityTestSlow(n int) (primeNumbers []int) {
	for i := 2; i <= n; i++ {
		for j := 2; j <= i; j++ {
			if i%j != 0 {
				continue
			}
			if i != j {
				break
			}
			primeNumbers = append(primeNumbers, j)
		}
	}
	return
}

func primalityTestFast(n int) (primeNumbers []int) {
	for i := 2; i <= n; i++ {
		s := math.Sqrt(float64(i))
		f := int(math.Floor(s))
		primal := true
		for j := 2; j <= f; j++ {
			if i%j == 0 {
				primal = false
				break
			}
		}
		if primal {
			primeNumbers = append(primeNumbers, i)
		}
	}
	return
}

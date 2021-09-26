package pprof

import (
	"fmt"
	"os"
	"runtime/pprof"
	"testing"
)

func TestLookupGoroutinesStack(t *testing.T) {
	prof := pprof.Lookup("goroutine")
	fmt.Println("-----------------------")
	prof.WriteTo(os.Stdout, 1)
	fmt.Println("-----------------------")
	prof.WriteTo(os.Stdout, 2)
	fmt.Println("-----------------------")
}

func TestWriteHeapProfile(t *testing.T) {
	pprof.WriteHeapProfile(os.Stdout)
}

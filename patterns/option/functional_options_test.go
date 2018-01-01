package option

import (
	"fmt"
	"testing"
	"time"
)

type fooOptions struct {
	verbose bool
	timeout time.Duration
}

type FooOption func(o *fooOptions)

func WithVerbose() FooOption {
	return func(o *fooOptions) {
		o.verbose = true
	}
}

func WithTimeout(timeout time.Duration) FooOption {
	return func(o *fooOptions) {
		o.timeout = timeout
	}
}

func DoFoo(opts ...FooOption) {
	o := new(fooOptions)
	for _, f := range opts {
		f(o)
	}

	// do something

	fmt.Printf("%+v", o)
}

func TestDoFoo(t *testing.T) {
	DoFoo(WithVerbose(), WithTimeout(500*time.Millisecond))
}

package wait

import "time"

// Options for call function
type Options struct {
	Timeout time.Duration
}

type optionFunc func(*Options)

func extract(opts ...optionFunc) Options {
	o := new(Options)
	for _, f := range opts {
		f(o)
	}
	return *o
}

// WithTimeout set timeout option
func WithTimeout(timeout time.Duration) optionFunc {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

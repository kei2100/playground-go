package http

import (
	"context"
	"net/http"
	"time"
)

type PathBasedTimeoutRoundTripper struct {
	RoundTripper   http.RoundTripper
	PathTimeout    map[string]time.Duration
	DefaultTimeout time.Duration
}

func (r *PathBasedTimeoutRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := r.RoundTripper
	if w == nil {
		w = http.DefaultTransport
	}
	to, ok := r.PathTimeout[req.URL.Path]
	if !ok {
		to = r.DefaultTimeout
	}
	if to > 0 {
		ctx, cancel := context.WithTimeout(req.Context(), to)
		defer cancel()
		req = req.Clone(ctx)
	}
	return w.RoundTrip(req)
}

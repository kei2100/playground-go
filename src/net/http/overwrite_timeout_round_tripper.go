package http

import (
	"context"
	"net/http"
	"time"
)

type OverwriteTimeoutRoundTripper struct {
	RoundTripper   http.RoundTripper
	TimeoutPerPath map[string]time.Duration
}

func (r *OverwriteTimeoutRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := r.RoundTripper
	if w == nil {
		w = http.DefaultTransport
	}
	if to, ok := r.TimeoutPerPath[req.URL.Path]; ok {
		ctx, cancel := context.WithTimeout(req.Context(), to)
		defer cancel()
		req = req.Clone(ctx)
	}
	return w.RoundTrip(req)
}

package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRoundTripper_RequestTimeout(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	req, err := http.NewRequest("GET", srv.URL, nil)
	if err != nil {
		t.Fatalf("new request failure: %v", err)
	}
	ctx, can := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer can()

	_, err = http.DefaultTransport.RoundTrip(req.WithContext(ctx))
	if g, w := err, context.DeadlineExceeded; g != w {
		t.Errorf("err got %v, want %v", g, w)
	}
}

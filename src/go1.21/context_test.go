package go1_21

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestContextWithoutCancel(t *testing.T) {
	pctx, cancel := context.WithCancel(context.Background())
	cancel()
	cctx := context.WithoutCancel(pctx)
	select {
	case <-cctx.Done():
		t.Error("cctx canceled")
	case <-time.After(10 * time.Millisecond):
		// ok
	}
}

func TestContextAfterFunc(t *testing.T) {
	// create test server
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)
		io.WriteString(w, "ok")
	})
	wrap := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.AfterFunc(r.Context(), func() {
			w.WriteHeader(499) // client closed
		})
		fn(w, r)
	})
	// send request
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	rec := httptest.NewRecorder()
	wrap(rec, req)
	if g, w := rec.Code, 499; g != w {
		t.Errorf("\ngot :%v\nwant:%v", rec.Code, 499)
	}
}

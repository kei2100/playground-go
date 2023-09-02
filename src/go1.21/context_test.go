package go1_21

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
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
		time.Sleep(100 * time.Millisecond)
		io.WriteString(w, "ok")
	})
	wrap := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = &mutexResponseWriter{ResponseWriter: w}
		stop := context.AfterFunc(r.Context(), func() {
			w.WriteHeader(499) // client closed
		})
		defer stop()
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
	// send request 2
	req, _ = http.NewRequest("GET", "/", nil)
	rec = httptest.NewRecorder()
	wrap(rec, req)
	if g, w := rec.Code, 200; g != w {
		t.Errorf("\ngot :%v\nwant:%v", rec.Code, 499)
	}
}

type mutexResponseWriter struct {
	mu sync.Mutex
	http.ResponseWriter
}

func (m *mutexResponseWriter) Write(b []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.ResponseWriter.Write(b)
}

func (m *mutexResponseWriter) WriteHeader(statusCode int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ResponseWriter.WriteHeader(statusCode)
}

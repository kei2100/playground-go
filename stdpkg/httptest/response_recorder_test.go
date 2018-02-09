package httptest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func hello(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("X-HELLO", "hello")
	w.Write([]byte("hello"))
}

func TestRecorder(t *testing.T) {
	t.Parallel()

	helloreq := httptest.NewRequest("GET", "/hellp", nil)
	w := httptest.NewRecorder()
	hello(w, helloreq)

	if g, w := w.Code, 200; g != w {
		t.Errorf("Code got %v, want %v", g, w)
	}
	if g, w := w.Header().Get("X-HELLO"), "hello"; g != w {
		t.Errorf("X-HELLO header got %v, want %v", g, w)
	}
	if g, w := w.Body.String(), "hello"; g != w {
		t.Errorf("Body got %v, want %v", g, w)
	}
}

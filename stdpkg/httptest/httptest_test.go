package httptest

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func echo() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		io.Copy(w, r.Body)
	}
}

func TestSimple(t *testing.T) {
	t.Parallel()

	sv := httptest.NewServer(http.HandlerFunc(echo()))
	defer sv.Close()

	r, err := http.Post(sv.URL, "text/plain", strings.NewReader("hello"))
	if err != nil {
		t.Errorf("http post got %v, want not error", err)
	}

	if g, w := r.StatusCode, 200; g != w {
		t.Errorf("StatusCode got %v, want %v", g, w)
	}

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		t.Fatal(err)
	}
	if g, w := string(b), "hello"; g != w {
		t.Errorf("response body got %v, want %v", g, w)
	}
}

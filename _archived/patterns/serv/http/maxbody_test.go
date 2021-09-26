package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kei2100/playground-go/patterns/serv/http/internal/response"
)

func TestMaxBody(t *testing.T) {
	s := Server{
		MaxBodyBytes: 3,
	}
	s.Route()
	s.router.Post("/echo", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			response.SendInternalServerError(w)
			return
		}
		defer r.Body.Close()
		w.Write(b)

	})

	tt := []struct {
		payload            string
		forceContentLength int64
		wantStatus         int
	}{
		{payload: "tes", wantStatus: 200},
		{payload: "test", wantStatus: 413},
		{payload: "tes", forceContentLength: 3, wantStatus: 200},
		{payload: "test", forceContentLength: 3, wantStatus: 500},
	}
	for i, te := range tt {
		r, _ := http.NewRequest("POST", "/echo", strings.NewReader(te.payload))
		if te.forceContentLength > 0 {
			r.ContentLength = te.forceContentLength
		}
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, r)

		if g, w := rec.Code, te.wantStatus; g != w {
			t.Errorf("status code got %v, want %v, at %d", g, w, i)
		}
		if te.wantStatus != 200 {
			continue
		}
		if g, w := rec.Body.String(), te.payload; g != w {
			t.Errorf("payload got %v, want %v, at %v", g, w, i)
		}
	}
}

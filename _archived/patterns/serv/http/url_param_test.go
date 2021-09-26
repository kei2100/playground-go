package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_withURLParam(t *testing.T) {
	s := Server{}
	s.Route()
	s.router.Post("/param/{param}", withURLParam(s.handleParam()))

	r, _ := http.NewRequest("POST", "/param/foo", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, r)

	if g, w := rec.Body.String(), "foo"; g != w {
		t.Errorf("body got %v, want %v", g, w)
	}
}

func Test_withURLParam2(t *testing.T) {
	s := Server{}
	s.Route()
	s.router.Post("/param/{zzz}/{yyy}", withURLParam2(s.handleParam2()))

	r, _ := http.NewRequest("POST", "/param/foo/bar", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, r)

	if g, w := rec.Body.String(), "foobar"; g != w {
		t.Errorf("body got %v, want %v", g, w)
	}
}

func (s *Server) handleParam() ParamHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, param string) {
		w.Write([]byte(param))
	}
}

func (s *Server) handleParam2() Param2HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, param1, param2 string) {
		w.Write([]byte(param1 + param2))
	}
}

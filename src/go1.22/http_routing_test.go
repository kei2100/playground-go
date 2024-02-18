package go1_22

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPRouting(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /foo/{value}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.PathValue("value")))
	})
	// send
	req, _ := http.NewRequest("POST", "/foo/bar", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "bar", rec.Body.String())
	// want 404
	req, _ = http.NewRequest("GET", "/foo/bar", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

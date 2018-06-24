package http

import (
	"bytes"
	"net/http"
	"time"
	"testing"
	"net/http/httptest"
	"io/ioutil"
	"reflect"
)

type downloader struct {
	*http.ServeMux
}

func NewDownloader() *downloader {
	mux := http.NewServeMux()
	mux.HandleFunc("/dl", func(w http.ResponseWriter, r *http.Request) {
		modtime := time.Now()
		content := bytes.NewReader([]byte("test content"))

		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Content-Disposition", `attachment; filename="content.txt"`)
		http.ServeContent(w, r, "content.txt", modtime, content)
	})
	return &downloader{mux}
}

func TestDownloader(t *testing.T) {
	t.Parallel()

	sv := httptest.NewServer(NewDownloader())
	defer sv.Close()

	res, err := http.Get(sv.URL + "/dl")
	if err != nil {
		t.Fatal(err)
	}
	b := res.Body
	defer b.Close()

	if !reflect.DeepEqual(res.Header["Content-Disposition"], []string{`attachment; filename="content.txt"`}) {
		t.Errorf("content-disposition got %v, want %v", res.Header["Content-Disposition"], []string{`attachment; filename="content.txt"`})
	}

	raw, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if g, w := string(raw), "test content"; g != w {
		t.Errorf("response body got %v, want %v", g, w)
	}
}

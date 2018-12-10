package handler_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	httplib "github.com/kei2100/playground-go/patterns/serv/http"
)

func servAndRecord(t *testing.T, method, path string, header http.Header, body io.Reader) httptest.ResponseRecorder {
	t.Helper()
	r, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Errorf("failed to create a request: %v", err)
		return httptest.ResponseRecorder{}
	}
	for h, vv := range header {
		for _, v := range vv {
			r.Header.Add(h, v)
		}
	}
	srv := httplib.Server{}
	srv.Route()
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, r)
	return *rec
}

func assertResponseCode(t *testing.T, gotCode, wantCode int) {
	t.Helper()
	if g, w := gotCode, wantCode; g != w {
		t.Errorf("response code got %v, want %v", g, w)
	}
}

func assertResponseJSON(t *testing.T, gotBody io.Reader, wantBody interface{}) {
	t.Helper()

	var wantJSON []byte
	rv := reflect.Indirect(reflect.ValueOf(wantBody))
	switch rv.Kind() {
	case reflect.String:
		wantJSON = []byte(rv.String())
	case reflect.Struct, reflect.Map:
		var err error
		wantJSON, err = json.Marshal(wantBody)
		if err != nil {
			t.Errorf("failed to marshal wantBody to json: %v", err)
			return
		}
	}

	got := make(map[string]interface{})
	dec := json.NewDecoder(gotBody)
	if err := dec.Decode(&got); err != nil {
		t.Errorf("failed to unmarshal got body to json: %v", err)
		return
	}
	want := make(map[string]interface{})
	if err := json.Unmarshal(wantJSON, &want); err != nil {
		t.Errorf("failed to unmarshal want body to json: %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("response json\ngot  %+v\nwant %+v", got, want)
	}
}

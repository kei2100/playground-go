package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	httplib "github.com/kei2100/playground-go/patterns/serv/http"
)

func sendRequest(t *testing.T, method, path string, body io.Reader) httptest.ResponseRecorder {
	t.Helper()
	r, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Errorf("failed to create a request: %v", err)
		return httptest.ResponseRecorder{}
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

func assertResponseJSON(t *testing.T, gotBody *bytes.Buffer, wantBody interface{}) {
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
	want := make(map[string]interface{})
	err := json.Unmarshal(gotBody.Bytes(), &got)
	if err != nil {
		t.Errorf("failed to unmarshal got body to json: %v", err)
		return
	}
	err = json.Unmarshal(wantJSON, &want)
	if err != nil {
		t.Errorf("failed to unmarshal want body to json: %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("response json\ngot  %+v\nwant %+v", got, want)
	}
}

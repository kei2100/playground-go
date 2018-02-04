package stub

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GreaterServer interface {
	Hello(http.ResponseWriter, *http.Request)
	GoodBy(http.ResponseWriter, *http.Request)
}

// Interfaceを埋め込んでスタブにしてしまう
type MinimalStubServer struct {
	GreaterServer
}

// 最低限のテスト対象メソッドのみ再定義する
func (s *MinimalStubServer) Hello(w http.ResponseWriter, r *http.Request) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	bw.WriteString("Hello!")
}

func TestStubInterface(t *testing.T) {
	s := new(MinimalStubServer)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.Hello)

	sv := httptest.NewServer(mux)
	defer sv.Close()

	r, err := http.Get(sv.URL + "/hello")
	if err != nil {
		t.Fatal(err)
	}
	if g, w := r.StatusCode, 200; g != w {
		t.Errorf("StatusCode got %v, want %v", g, w)
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
	}
	if g, w := string(b), "Hello!"; g != w {
		t.Errorf("Body got %v, want %v", g, w)
	}
}

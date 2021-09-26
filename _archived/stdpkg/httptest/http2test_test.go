package httptest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

	"fmt"

	"golang.org/x/net/http2"
)

func TestHTTP2ServerClient(t *testing.T) {
	t.Parallel()

	srv, client := http2StartedServerClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()

	res, err := client.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	t.Logf("protocol ver: %v", res.Proto)
	if !res.ProtoAtLeast(2, 0) {
		t.Fatal("unexpected protocol ver")
	}
}

func TestHTTP2ServerPush(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("testdata/images"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if pusher, ok := w.(http.Pusher); ok {
			// Developer toolのNetworkタブで確認すると、Initiatorが「Push/(index)」になっているのが確認できる。
			pusher.Push("/images/music.png", nil)
		}
		w.Header().Add("Content-Type", "text/html")
		fmt.Fprintf(w, `<html><body><img src="/images/music.png"></body></html>`)
	})

	srv, _ := http2StartedServerClient(t, mux)
	fmt.Printf("server listening on %v\n", srv.Listener.Addr().String())
	defer srv.Close()

	t.Skip("skip for auto testing")
	time.Sleep(time.Hour)
}

func http2StartedServerClient(t *testing.T, handler http.Handler) (*httptest.Server, *http.Client) {
	t.Helper()

	srv := httptest.NewUnstartedServer(handler)
	// see https://github.com/golang/go/blob/bf9f1c15035ab9bb695a9a3504e465a1896b4b8c/src/net/http/server.go#L3065
	h2ServerConfig := &http2.Server{
		NewWriteScheduler: func() http2.WriteScheduler { return http2.NewPriorityWriteScheduler(nil) },
	}
	if err := http2.ConfigureServer(srv.Config, h2ServerConfig); err != nil {
		t.Fatal(err)
	}
	srv.TLS = srv.Config.TLSConfig
	srv.StartTLS()

	client := srv.Client()
	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("unexpected client transport type %T", client.Transport)
	}
	if err := http2.ConfigureTransport(tr); err != nil {
		t.Fatal(err)
	}

	return srv, client
}

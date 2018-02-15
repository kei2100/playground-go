package httptest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/http2"
)

func TestHTTP2(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	// TODO register routes

	srv := httptest.NewUnstartedServer(mux)
	setNextProtoDefaults(srv.Config)
}

// setNextProtoDefaults configures HTTP/2.
// see https://github.com/golang/go/blob/bf9f1c15035ab9bb695a9a3504e465a1896b4b8c/src/net/http/server.go#L3065
func setNextProtoDefaults(srv *http.Server) {
	if srv.TLSNextProto != nil {
		return
	}
	conf := &http2.Server{
		NewWriteScheduler: func() http2.WriteScheduler { return http2.NewPriorityWriteScheduler(nil) },
	}
	http2.ConfigureServer(srv, conf)
}

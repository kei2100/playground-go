package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTest(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	t.Cleanup(svr.Close)
	// create trace
	trace := httptrace.ClientTrace{
		GetConn: nil,
		GotConn: func(i httptrace.GotConnInfo) {
			if i.WasIdle {
				fmt.Printf("got idle conn. idle time %s\n", i.IdleTime)
			} else {
				fmt.Println("got new conn.")
			}
		},
		PutIdleConn:          nil,
		GotFirstResponseByte: nil,
		Got100Continue:       nil,
		Got1xxResponse:       nil,
		DNSStart:             nil,
		DNSDone:              nil,
		ConnectStart:         nil,
		ConnectDone:          nil,
		TLSHandshakeStart:    nil,
		TLSHandshakeDone:     nil,
		WroteHeaderField:     nil,
		WroteHeaders:         nil,
		Wait100Continue:      nil,
		WroteRequest:         nil,
	}
	orig, _ := http.NewRequest("GET", svr.URL, nil)
	// send request
	req := orig.WithContext(httptrace.WithClientTrace(context.Background(), &trace))
	_, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	// re-send request
	req = orig.WithContext(httptrace.WithClientTrace(context.Background(), &trace))
	_, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
}

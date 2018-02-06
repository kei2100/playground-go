package main

import (
	"log"
	"time"

	"net"
	"net/http"
	"strconv"

	"github.com/google/gops/agent"
)

func main() {
	// Options.NoShutdownCleanup = false(デフォルト)だとSIGINTをトラップして
	// os.Exit(1)する模様
	if err := agent.Listen(&agent.Options{}); err != nil {
		log.Fatal(err)
	}
	defer agent.Close()
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("http server listening on %v", ln.Addr().String())
	http.Serve(ln, serveMux())
}

func serveMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/delay", delay)
	return mux
}

func index(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
}

func delay(w http.ResponseWriter, r *http.Request) {
	sec := r.URL.Query()["sec"]
	var s string
	if len(sec) > 0 {
		s = sec[0]
	}
	d, err := strconv.Atoi(s)
	if err != nil {
		d = 0
	}
	if d > 300 {
		d = 300
	}
	time.Sleep(time.Duration(d) * time.Second)
	w.WriteHeader(200)
}

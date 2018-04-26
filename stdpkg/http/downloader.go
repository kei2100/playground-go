package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	fmt.Println("listening on " + l.Addr().String())

	http.HandleFunc("/dl", func(w http.ResponseWriter, r *http.Request) {
		modtime := time.Now()
		content := bytes.NewReader([]byte("test content"))

		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Content-Disposition", `attachment; filename="content.txt"`)
		http.ServeContent(w, r, "content.txt", modtime, content)
	})
	http.Serve(l, http.DefaultServeMux)
}

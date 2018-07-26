package main

import (
	"net/http"

	"fmt"
	"net"

	"log"

	"github.com/kei2100/playground-go/util/http/proxy/fxy/proxy"
)

func main() {
	cfg := proxy.Config{
		URLConfig: proxy.URLConfig{
			Server: "https://www.google.com",
		},
	}

	ln, err := net.Listen("tcp", "localhost:18888")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Printf("listening on %s", ln.Addr())

	if err := cfg.Load(); err != nil {
		log.Fatalln(err)
	}
	err = http.Serve(ln, &proxy.Server{Config: cfg})
	fmt.Println(err)
}

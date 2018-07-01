package main

import (
	"fmt"
	"net"
	"net/http"

	"net/url"

	"github.com/kei2100/playground-go/util/http/proxy/fxy/proxy"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	fmt.Printf("listening on %s", ln.Addr())

	dest, err := url.Parse("https://www.google.com")
	if err != nil {
		panic(err)
	}
	sv := proxy.Server{Destination: dest}
	http.Serve(ln, &sv)
}

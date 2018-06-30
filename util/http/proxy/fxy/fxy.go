package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/kei2100/playground-go/util/http/proxy/fxy/proxy"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	fmt.Printf("listening on %s", ln.Addr())
	http.Serve(ln, proxy.New())
}

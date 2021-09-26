package gateway_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/kei2100/playground-go/util/tls/gateway"
)

func ExampleListen() {
	remote, received := startRemoteServer()
	defer remote.Close()

	gtw, err := gateway.Listen()
	if err != nil {
		log.Fatalf("gateway_test: failed to listen gateway")
	}
	defer gtw.Close()

	router := gateway.NewRouter(remote.Addr().String(), gateway.WithNoTLS())
	go gtw.Serve(router)

	sendMessage(gtw.Addr(), []byte("hello"))

	select {
	case rmsg := <-received:
		// Output:
		// hello
		//
		fmt.Println(string(rmsg))
	case <-time.After(100 * time.Millisecond):
		log.Fatalln("gateway_test: timeout exceeded while waiting for receive the message")
	}
}

func startRemoteServer() (net.Listener, <-chan []byte) {
	received := make(chan []byte)

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatalf("gateway_test: failed to listen remote server: %v", err)
	}

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("gateway_test: failed to accept incoming connection: %v", err)
		}
		defer conn.Close()

		b, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Fatalf("gateway_test: failed to read bytes: %v", err)
		}

		received <- b
		close(received)
	}()

	return ln, received
}

func sendMessage(addr net.Addr, msg []byte) {
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		log.Fatalf("gateway_test: faield to dial to %v: %v", addr.String(), err)
	}
	defer conn.Close()
	if _, err := conn.Write(msg); err != nil {
		log.Fatalf("gateway_test: faield to write msg: %v", err)
	}
}

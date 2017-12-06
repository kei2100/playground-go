package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"
	"sync"
)

var destAddr = "www.google.co.jp:80"
var localAddr = "localhost:3100"
var useTLS = false

func init() {
	flag.StringVar(&destAddr, "d", "www.google.co.jp:80", "destination address")
	flag.StringVar(&localAddr, "l", "localhost:3100", "local address")
	flag.BoolVar(&useTLS, "s", false, "use tls")
}

func main() {
	flag.Parse()

	addr, err := net.ResolveTCPAddr("tcp", localAddr)
	if err != nil {
		log.Fatalf("failed to resolve addr:%v", err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("failed to accept the connection:%v", err)
			continue
		}
		go handle(conn)
	}
}

func handle(src net.Conn) {
	var dest net.Conn
	var err error

	if useTLS {
		dest, err = tls.Dial("tcp", destAddr, &tls.Config{})
	} else {
		dest, err = net.Dial("tcp", destAddr)
	}
	if err != nil {
		log.Fatalf("dial failed %v", err)
	}

	var wg sync.WaitGroup

	// FIXME 現状だと片側クローズされるとハーフオープン状態になる

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer src.Close()
		if _, err := io.Copy(dest, src); err != nil {
			log.Printf("failed to io copy src to destination: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer dest.Close()
		if _, err := io.Copy(src, dest); err != nil {
			log.Printf("failed to io copy destination to src: %v", err)
		}
	}()

	wg.Wait()
	log.Println("handle done")
}

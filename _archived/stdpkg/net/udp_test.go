package net

import (
	"fmt"
	"log"
	"net"
	"testing"
)

type udpEchoServer struct {
	conn net.PacketConn
}

func listenUDP() (*udpEchoServer, error) {
	conn, err := net.ListenPacket("udp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("net: failed to listen: %v", err)
	}

	s := udpEchoServer{conn: conn}
	return &s, nil
}

func (s *udpEchoServer) serve() error {
	for {
		buf := make([]byte, 1500)
		length, raddr, err := s.conn.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("net: failed to read from packet connection: %v", err)
		}

		if _, err := s.conn.WriteTo(buf[:length], raddr); err != nil {
			return fmt.Errorf("net: failed to write to packet connection: %v", err)
		}
	}
}

func (s *udpEchoServer) close() error {
	return s.conn.Close()
}

func TestUDPEcho(t *testing.T) {
	srv, err := listenUDP()
	if err != nil {
		t.Fatal(err)
	}
	defer srv.close()
	go func() {
		log.Printf("net: echo server stopped: %v", srv.serve())
	}()

	conn, err := net.Dial("udp4", srv.conn.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1500)
	length, err := conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if g, w := string(buf[:length]), "hello"; g != w {
		t.Errorf(" got %v, want %v", g, w)
	}
}

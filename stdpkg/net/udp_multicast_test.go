package net

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

type multicastUnixTimeServer struct {
	stopSig chan struct{}
	done    chan struct{}
}

func (s *multicastUnixTimeServer) serve() error {
	conn, err := net.Dial("udp4", "224.0.0.1:9999")
	if err != nil {
		return fmt.Errorf("net: failed to dial: %v", err)
	}
	defer conn.Close()

	s.stopSig, s.done = make(chan struct{}), make(chan struct{})
	b := make([]byte, 8)

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			binary.LittleEndian.PutUint64(b, uint64(time.Now().Unix()))
			if _, err := conn.Write(b); err != nil {
				return fmt.Errorf("net: failed to write to the connection: %v", err)
			}
		case <-s.stopSig:
			close(s.done)
			return fmt.Errorf("net: serve stopped normaly")
		}
	}
}

func (s *multicastUnixTimeServer) stop() {
	s.stopSig <- struct{}{}
	<-s.done
}

func TestUDPMulticast(t *testing.T) {
	srv := &multicastUnixTimeServer{}
	go func() {
		log.Println(srv.serve())
	}()
	defer srv.stop()

	addr, err := net.ResolveUDPAddr("udp4", "224.0.0.1:9999")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	cnt := 0

	for {
		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout exceeded")
		default:
			if cnt > 2 {
				return
			}

			b := make([]byte, 8)
			l, _, err := conn.ReadFrom(b)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(binary.LittleEndian.Uint64(b[:l]))
			cnt++
		}
	}
}

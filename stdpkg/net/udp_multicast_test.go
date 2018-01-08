package net

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"
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
	close(s.stopSig)
	<-s.done
}

type multicastClient struct {
	stopSig chan struct{}
	done    chan struct{}
}

func (c *multicastClient) listen() (<-chan []byte, error) {
	addr, err := net.ResolveUDPAddr("udp4", "224.0.0.1:9999")
	if err != nil {
		return nil, fmt.Errorf("net: failed to resolve addr: %v", err)
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("net: failed to listen: %v", err)
	}

	data := make(chan []byte)
	c.stopSig, c.done = make(chan struct{}), make(chan struct{})

	go func() {
		defer func() {
			conn.Close()
			close(data)
			close(c.done)
		}()

		for {
			select {
			case <-c.stopSig:
				return
			default:
				b := make([]byte, 8)
				l, _, err := conn.ReadFromUDP(b)
				if err != nil {
					fmt.Printf("net: faield to read from the connection: %v", err)
					return
				}
				data <- b[:l]
			}
		}
	}()

	return data, nil
}

func (c *multicastClient) stop() {
	close(c.stopSig)
	<-c.done
}

func TestUDPMulticast(t *testing.T) {
	srv := &multicastUnixTimeServer{}
	go func() {
		fmt.Println(srv.serve())
	}()
	defer srv.stop()

	c1, c2 := &multicastClient{}, &multicastClient{}

	d1, err := c1.listen()
	if err != nil {
		t.Fatal(err)
	}
	defer c1.stop()

	d2, err := c2.listen()
	if err != nil {
		t.Fatal(err)
	}
	defer c2.stop()

	var cnt int32 = 0
	ctx, can := context.WithTimeout(context.Background(), 5*time.Second)
	defer can()

	go func() {
		for d := range d1 {
			atomic.AddInt32(&cnt, 1)
			fmt.Printf("d1 %v\n", binary.LittleEndian.Uint64(d))
		}
	}()

	go func() {
		for d := range d2 {
			atomic.AddInt32(&cnt, 1)
			fmt.Printf("d2 %v\n", binary.LittleEndian.Uint64(d))
		}
	}()

	for {
		select {
		case <-ctx.Done():
			t.Fatal("timeout exceeded")
		default:
			if cnt > 5 {
				return
			}
		}
	}
}

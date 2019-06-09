package mq

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"testing"
)

func TestMQ_Receive(t *testing.T) {
	t.Run("Context canceled", func(t *testing.T) {
		m, err := Create(randHex(8), os.O_CREATE, 0600)
		if err != nil {
			t.Fatal(err)
		}
		defer m.Close()

		ctx, cancel := context.WithCancel(context.Background())
		go func() { cancel() }()

		p := make([]byte, 8)
		n, err := m.ReceiveContext(ctx, p)

		if err == nil || err != context.Canceled {
			t.Errorf("got %v, want an error", err)
		}
		if g, w := n, 0; g != w {
			t.Errorf("n got %v, want %v", g, w)
		}
	})

	t.Run("Closed MQ", func(t *testing.T) {
		// TODO
	})
}

func randHex(n int) string {
	p := make([]byte, n)
	rand.Read(p)
	return fmt.Sprintf("%x", p)
}

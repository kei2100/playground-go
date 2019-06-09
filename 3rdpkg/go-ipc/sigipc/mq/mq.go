package mq

import (
	"context"
	"fmt"
	"os"

	"bitbucket.org/avd/go-ipc/mq"
	"golang.org/x/sync/errgroup"
)

// MQ is a message queue
// TODO
type MQ struct {
	closed chan struct{}
	name   string
	mq     mq.Messenger
}

// Create MQ
//	flag - create flags. You can specify:
//		os.O_EXCL
//		mq.O_NONBLOCK
func Create(name string, flag int, perm os.FileMode) (*MQ, error) {
	mmq, err := mq.New(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return &MQ{closed: make(chan struct{}), name: name, mq: mmq}, nil
}

// Open MQ
//	flag - create flags. You can specify:
//		0 or mq.O_NONBLOCK
func Open(name string, flag int) (*MQ, error) {
	mmq, err := mq.Open(name, flag)
	if err != nil {
		return nil, err
	}
	return &MQ{closed: make(chan struct{}), name: name, mq: mmq}, nil
}

// Remove MQ
func Remove(name string) error {
	return mq.Destroy(name)
}

// Close MQ
func (q *MQ) Close() error {
	close(q.closed)
	return q.mq.Close()
}

// Receive msg
func (q *MQ) Receive(p []byte) (int, error) {
	return q.mq.Receive(p)
}

// ReceiveContext msg
func (q *MQ) ReceiveContext(ctx context.Context, p []byte) (int, error) {
	mmq, err := mq.Open(q.name, mq.O_NONBLOCK)
	if err != nil {
		return 0, fmt.Errorf("mq: failed to open the non-block mq: %+v", err)
	}
	defer mmq.Close()

	done := make(chan struct{})
	var eg errgroup.Group
	eg.Go(func() error {
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			if err := mmq.Send([]byte{}); err != nil {
				fmt.Printf("mq: failed to send: %+v", err) // TODO
			}
			return ctx.Err()
		case <-q.closed:
			if err := mmq.Send([]byte{}); err != nil {
				fmt.Printf("mq: failed to send: %+v", err) // TODO
			}
			return fmt.Errorf("mq: closed")
		}
	})

	var n int
	eg.Go(func() error {
		defer close(done)
		var err error
		n, err = q.mq.Receive(p)
		if err != nil {
			return fmt.Errorf("mq: failed to receive: %+v", err)
		}
		return nil
	})

	return n, eg.Wait()
}

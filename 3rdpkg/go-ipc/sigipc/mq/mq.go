package mq

import (
	"bitbucket.org/avd/go-ipc/mq"
	"context"
	"fmt"
	"os"
)

// MQ is a message queue
// TODO
type MQ struct {
	name string
	mq mq.Messenger
}

// Open MQ
func Open(name string) (*MQ, error) {
	mmq, err := mq.Open(name, os.O_CREATE)
	if err != nil {
		return nil, fmt.Errorf("mq: failed to open the message queue: %+v", err)
	}
	return &MQ{name: name, mq: mmq}, nil
}

// Close MQ
func (q *MQ) Close() error {
	err := mq.Destroy(q.name)
	if err != nil {
		return fmt.Errorf("mq: failed to destroy the message queue: %+v", err)
	}
	return nil
}

// Receive msg
func (q *MQ) Receive(p []byte, ctx context.Context) (int, error) {
	// duplex	
	return 0, nil
}
package wait

import (
	"context"
	"fmt"
)

// ReceiveStuct waits receive from the channel
func ReceiveStuct(ch <-chan struct{}, opts ...optionFunc) error {
	o := extract(opts...)

	if o.Timeout == 0 {
		return ReceiveStructContext(context.Background(), ch)
	}
	ctx, can := context.WithTimeout(context.Background(), o.Timeout)
	defer can()
	return ReceiveStructContext(ctx, ch)
}

// ReceiveStructContext waits receive from the channel with context
func ReceiveStructContext(ctx context.Context, ch <-chan struct{}) error {
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait: %v", ctx.Err())
	}
}

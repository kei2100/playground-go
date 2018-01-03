package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// SomeTask struct
type SomeTask struct {
	Wait time.Duration
}

// Do something with context
func (t *SomeTask) Do(ctx context.Context) error {
	var success error = nil
	defer t.free()

	select {
	case <-time.After(t.Wait):
		return success
	case <-ctx.Done():
		return fmt.Errorf("context canceled :%v", ctx.Err())
	}
}

// Free something resource
func (t *SomeTask) free() { fmt.Println("free something resource") }

func TestContextWithTimeout(t *testing.T) {
	tests := []struct {
		subject   string
		taskWait  time.Duration
		timeout   time.Duration
		wantError bool
	}{
		{
			subject:   "not timeout",
			taskWait:  1 * time.Millisecond,
			timeout:   10 * time.Millisecond,
			wantError: false,
		},
		{
			subject:   "timeout",
			taskWait:  10 * time.Millisecond,
			timeout:   1 * time.Millisecond,
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			task := &SomeTask{Wait:tt.taskWait}
			ctx, can := context.WithTimeout(context.Background(), tt.timeout)
			defer can()
			err := task.Do(ctx)

			if tt.wantError && err == nil {
				t.Error("want error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("want nil, got %v", err)
			}
		})
	}
}

package time

import (
	"fmt"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(500 * time.Millisecond)
	cnt := 0
	for {
		select {
		case <-ticker.C:
			if cnt > 2 {
				return
			}
			fmt.Printf("tick %v\n", cnt)
			cnt++
		case <-time.After(5 * time.Second):
			t.Fatal("timeout exceeded")
		}
	}
}

package go1_24

import (
	"fmt"
	"testing"
)

func TestTestContext(t *testing.T) {
	t.Cleanup(func() {
		fmt.Println("TestTestContext cleanup start")

		select {
		case <-t.Context().Done():
			// Cleanup では t.Context().Done している状態になるのでこちらに入る
			fmt.Println("TestTestContext done")
		default:
			fmt.Println("TestTestContext not done")
		}

		fmt.Println("TestTestContext cleanup end")
	})

	fmt.Println("TestTestContext start")
	select {
	case <-t.Context().Done():
		fmt.Println("TestTestContext done")
	default:
		fmt.Println("TestTestContext not done")
	}
	fmt.Println("TestTestContext end")
}

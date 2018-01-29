package pipeline

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"golang.org/x/text/width"
)

func TestPipeline(t *testing.T) {
	ctx, can := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer can()

	c := gen(ctx, "ｈｅｌｌｏ", "ｗｏｒｌｄ", "ｇｏｌａｎｇ")
	rs := toUpper(ctx, toNarrow(ctx, c))

	var strs []string
	for r := range rs {
		strs = append(strs, r)
		// このsleepを付けるとctxがタイムアウトする
		//time.Sleep(110 * time.Millisecond)
	}
	fmt.Println(strings.Join(strs, " "))
}

// generate pipeline
func gen(ctx context.Context, strs ...string) <-chan string {
	out := make(chan string)

	go func() {
		defer func() {
			log.Println("gen: finished. close out channel")
			close(out)
		}()

		for _, s := range strs {
			select {
			case out <- s:
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					log.Printf("gen: context was canceled or deadline exceeded: %v", err)
				}
				return
			}
		}
	}()

	return out
}

func toNarrow(ctx context.Context, in <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer func() {
			log.Println("toNarrow: finished. close out channel")
			close(out)
		}()

		for s := range in {
			select {
			case out <- width.Narrow.String(s):
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					log.Printf("toNarrow: context was canceled or deadline exceeded: %v", err)
				}
				return
			}
		}
	}()

	return out
}

func toUpper(ctx context.Context, in <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer func() {
			log.Println("toUpper: finished. close out channel")
			close(out)
		}()

		for s := range in {
			select {
			case out <- strings.ToUpper(s):
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					log.Printf("toUpper: context was canceled or deadline exceeded: %v", err)
				}
				return
			}
		}
	}()

	return out
}

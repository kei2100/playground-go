package concurrency

import (
	"testing"

	"fmt"
	"strings"

	"golang.org/x/text/width"
)

func TestPipeline(t *testing.T) {
	t.Parallel()

	c := gen("ｈｅｌｌｏ", "ｗｏｒｌｄ")
	rs := toUpper(toNarrow(c))

	var strs []string
	for r := range rs {
		strs = append(strs, r)
	}
	fmt.Println(strings.Join(strs, " "))
}

// generate pipeline
func gen(strs ...string) <-chan string {
	out := make(chan string)
	go func() {
		for _, s := range strs {
			out <- s
		}
		close(out)
	}()
	return out
}

func toNarrow(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for s := range in {
			out <- width.Narrow.String(s)
		}
		close(out)
	}()
	return out
}

func toUpper(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for s := range in {
			out <- width.Narrow.String(s)
		}
		close(out)
	}()
	return out
}

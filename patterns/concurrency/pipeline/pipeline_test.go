package pipeline

import (
	gocontext "context"
	"fmt"
	"log"
	"strings"
	"testing"
	"unicode"

	"github.com/kei2100/playground-go/util/context"
	"golang.org/x/text/width"
)

func TestPipeline(t *testing.T) {
	ctx, can := context.Failure(gocontext.Background())
	defer can()

	c := gen(ctx, "ｈｅｌｌｏ", "ｗｏｒｌｄ", "ｇｏｌａｎｇ")
	//c := gen(ctx, "ｈｅｌｌｏ", "ｗｏｒｌｄ", "ｇｏｌａｎｇ", string(unicode.ReplacementChar))	// error pattern
	rs := toUpper(ctx, toNarrow(ctx, c))

	var ss []string

loop:
	for {
		select {
		case s, ok := <-rs:
			if !ok {
				break loop
			}
			ss = append(ss, s)
		case <-ctx.Done():
			fmt.Printf("pipeline failed: %v\n", ctx.Err())
			return
		}
	}

	fmt.Println(strings.Join(ss, " "))
}

// generate pipeline
func gen(ctx gocontext.Context, strs ...string) <-chan string {
	out := make(chan string)

	go func() {
		defer func() {
			log.Println("gen: finished. close out channel")
			close(out)
		}()

		for _, s := range strs {
			select {
			case out <- s:
				continue
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func toNarrow(ctx context.FailureContext, in <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer func() {
			log.Println("toNarrow: finished. close out channel")
			close(out)
		}()

		for s := range in {
			select {
			case out <- width.Narrow.String(s):
				continue
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func toUpper(ctx context.FailureContext, in <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer func() {
			log.Println("toUpper: finished. close out channel")
			close(out)
		}()

		for s := range in {
			s := strings.ToUpper(s)
			if strings.ContainsRune(s, unicode.ReplacementChar) {
				ctx.Fail(fmt.Errorf("invalid code points contains: %v", s))
				return
			}
			select {
			case out <- s:
				continue
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

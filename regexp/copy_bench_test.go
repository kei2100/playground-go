package regexp

import (
	"regexp"
	"testing"
)

func BenchmarkRegexpShared(b *testing.B) {
	x := []byte("this is a long line that contains foo bar baz")
	re := regexp.MustCompile("foo (ba+r)? baz")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			re.Match(x)
		}
	})
}

func BenchmarkRegexpCopied(b *testing.B) {
	x := []byte("this is a long line that contains foo bar baz")
	re := regexp.MustCompile("foo (ba+r)? baz")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		re := re.Copy()
		for pb.Next() {
			re.Match(x)
		}
	})
}

package regexp

import (
	"regexp"
	"testing"
)

func BenchmarkRegexpShared(b *testing.B) {
	x := "this is a long line that contains foo bar baz"
	re := regexp.MustCompile("foo (ba+r)? baz")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			re.FindStringSubmatchIndex(x)
		}
	})
}

func BenchmarkRegexpCopied(b *testing.B) {
	x := "this is a long line that contains foo bar baz"
	re := regexp.MustCompile("foo (ba+r)? baz")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		re := re.Copy()
		for pb.Next() {
			re.FindStringSubmatchIndex(x)
		}
	})
}

func BenchmarkRegexpFindStringSubmatchIndex(b *testing.B) {
	r := regexp.MustCompile(`\A(\d{1,2}):(\d{1,2})(?::(\d{1,2}))?`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.FindStringSubmatchIndex("12:11:10")
	}
}

func BenchmarkRegexpFindStringSubmatch(b *testing.B) {
	r := regexp.MustCompile(`\A(\d{1,2}):(\d{1,2})(?::(\d{1,2}))?`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.FindStringSubmatch("12:11:10")
	}
}

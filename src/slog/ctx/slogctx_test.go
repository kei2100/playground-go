package slogctx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"sync"
	"testing"
	"testing/slogtest"
)

func newTestHandler() (*Handler, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return NewHandler(slog.NewJSONHandler(buf, nil)), buf
}

func parseLines(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, line := range bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte{'\n'}) {
		if len(line) == 0 {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal(line, &m); err != nil {
			t.Fatalf("unmarshal %q: %v", line, err)
		}
		out = append(out, m)
	}
	return out
}

// TestHandlerComplies は slog.Handler 契約を slogtest で網羅検証する。
func TestHandlerComplies(t *testing.T) {
	h, buf := newTestHandler()
	results := func() []map[string]any { return parseLines(t, buf) }
	if err := slogtest.TestHandler(h, results); err != nil {
		t.Fatal(err)
	}
}

// TestInjectsCtxAttrAsTopLevel は WithAttr した属性がトップレベルに出ることを確認する。
func TestInjectsCtxAttrAsTopLevel(t *testing.T) {
	h, buf := newTestHandler()
	logger := slog.New(h)

	ctx := WithAttr(context.Background(), slog.String("trace_id", "abc123"))
	logger.InfoContext(ctx, "hello")

	ms := parseLines(t, buf)
	if len(ms) != 1 {
		t.Fatalf("want 1 line, got %d", len(ms))
	}
	if v := ms[0]["trace_id"]; v != "abc123" {
		t.Errorf("trace_id not at top level: %v", ms[0])
	}
}

// TestCtxAttrStaysTopLevelUnderWithGroup は WithGroup が積まれていても
// ctx 属性がトップレベルに出ることを確認する (このパッケージの主目的)。
func TestCtxAttrStaysTopLevelUnderWithGroup(t *testing.T) {
	h, buf := newTestHandler()
	logger := slog.New(h).WithGroup("g1").WithGroup("g2")

	ctx := WithAttr(context.Background(), slog.String("trace_id", "abc123"))
	logger.InfoContext(ctx, "hello", "k", "v")

	ms := parseLines(t, buf)
	if len(ms) != 1 {
		t.Fatalf("want 1 line, got %d", len(ms))
	}
	m := ms[0]

	// ctx 由来はトップレベル
	if v := m["trace_id"]; v != "abc123" {
		t.Errorf("trace_id not at top level: %v", m)
	}
	// record 由来は g1.g2 配下
	g1, ok := m["g1"].(map[string]any)
	if !ok {
		t.Fatalf("group g1 not found: %v", m)
	}
	g2, ok := g1["g2"].(map[string]any)
	if !ok {
		t.Fatalf("group g2 not found: %v", g1)
	}
	if v := g2["k"]; v != "v" {
		t.Errorf("k=v not in g1.g2: %v", g2)
	}
}

// TestCtxAttrChainInheritsAndOrders は親 ctx の属性が子 ctx でも見えること、
// 出力順がルート→子であることを確認する。
func TestCtxAttrChainInheritsAndOrders(t *testing.T) {
	parent := WithAttr(context.Background(), slog.String("a", "1"))
	child := WithAttr(parent, slog.String("b", "2"))

	got := attrsFromContext(child)
	if len(got) != 2 {
		t.Fatalf("want 2 attrs, got %d: %v", len(got), got)
	}
	if got[0].Key != "a" || got[1].Key != "b" {
		t.Errorf("order: want [a, b], got [%s, %s]", got[0].Key, got[1].Key)
	}
}

// TestWithAttrDoesNotMutateParentCtx は子 ctx に WithAttr しても、
// 同じ親から派生した別の子 ctx に影響しないことを確認する (リンクリストの不変性)。
func TestWithAttrDoesNotMutateParentCtx(t *testing.T) {
	parent := WithAttr(context.Background(), slog.String("a", "1"))

	childX := WithAttr(parent, slog.String("x", "X"))
	childY := WithAttr(parent, slog.String("y", "Y"))

	gotX := attrsFromContext(childX)
	gotY := attrsFromContext(childY)
	gotP := attrsFromContext(parent)

	if len(gotP) != 1 || gotP[0].Key != "a" {
		t.Errorf("parent should only have a: %v", gotP)
	}
	if len(gotX) != 2 || gotX[1].Key != "x" {
		t.Errorf("childX should be [a, x]: %v", gotX)
	}
	if len(gotY) != 2 || gotY[1].Key != "y" {
		t.Errorf("childY should be [a, y]: %v", gotY)
	}
}

// TestWithGroupForkSafety は同じ base から WithGroup を複数派生しても
// 互いに汚染しないことを確認する (slices.Clone の効果検証)。
func TestWithGroupForkSafety(t *testing.T) {
	h, buf := newTestHandler()
	base := slog.New(h).WithGroup("a").WithGroup("b").WithGroup("c")

	base.WithGroup("d").Info("ma", "k", "vA")
	base.WithGroup("e").Info("mb", "k", "vB")

	ms := parseLines(t, buf)
	if len(ms) != 2 {
		t.Fatalf("want 2 lines, got %d", len(ms))
	}

	deep := func(m map[string]any, keys ...string) any {
		var cur any = m
		for _, k := range keys {
			mm, ok := cur.(map[string]any)
			if !ok {
				return nil
			}
			cur = mm[k]
		}
		return cur
	}

	if v := deep(ms[0], "a", "b", "c", "d", "k"); v != "vA" {
		t.Errorf("line0: a.b.c.d.k = %v, want vA. full=%v", v, ms[0])
	}
	if v := deep(ms[1], "a", "b", "c", "e", "k"); v != "vB" {
		t.Errorf("line1: a.b.c.e.k = %v, want vB. full=%v", v, ms[1])
	}
}

// TestWithAttrsInGroupForkSafety は WithGroup 後に WithAttrs を複数派生しても
// 互いに汚染しないことを確認する (groupState 値スライス化の効果検証)。
func TestWithAttrsInGroupForkSafety(t *testing.T) {
	h, buf := newTestHandler()
	base := slog.New(h).WithGroup("g")

	base.With("x", 1).Info("ma")
	base.With("y", 2).Info("mb")

	ms := parseLines(t, buf)
	if len(ms) != 2 {
		t.Fatalf("want 2 lines, got %d", len(ms))
	}

	gA, _ := ms[0]["g"].(map[string]any)
	gB, _ := ms[1]["g"].(map[string]any)

	if _, ok := gA["x"]; !ok {
		t.Errorf("childA missing x: %v", ms[0])
	}
	if _, ok := gA["y"]; ok {
		t.Errorf("childA contaminated with y: %v", ms[0])
	}
	if _, ok := gB["y"]; !ok {
		t.Errorf("childB missing y: %v", ms[1])
	}
	if _, ok := gB["x"]; ok {
		t.Errorf("childB contaminated with x: %v", ms[1])
	}
}

// TestEmptyGroupAndAttrs は空 group / 空 attrs が無視されることを確認する。
func TestEmptyGroupAndAttrs(t *testing.T) {
	h := NewHandler(slog.NewJSONHandler(&bytes.Buffer{}, nil))
	if got := h.WithGroup(""); got != h {
		t.Errorf("WithGroup(\"\") should return self")
	}
	if got := h.WithAttrs(nil); got != h {
		t.Errorf("WithAttrs(nil) should return self")
	}
}

// TestWithAttrConcurrent は同じ親 ctx に対し並行に WithAttr を呼んでも
// データレースしないこと、各 ctx が独立した結果を返すことを確認する。
// go test -race で実行すること。
func TestWithAttrConcurrent(t *testing.T) {
	parent := WithAttr(context.Background(), slog.String("root", "yes"))

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx := WithAttr(parent, slog.Int("i", i))
			attrs := attrsFromContext(ctx)
			if len(attrs) != 2 {
				t.Errorf("goroutine %d: want 2 attrs, got %d: %v", i, len(attrs), attrs)
				return
			}
			if attrs[0].Key != "root" || attrs[1].Key != "i" {
				t.Errorf("goroutine %d: order wrong: %v", i, attrs)
			}
			if int(attrs[1].Value.Int64()) != i {
				t.Errorf("goroutine %d: i=%d, got %v", i, i, attrs[1])
			}
		}(i)
	}
	wg.Wait()
}

// BenchmarkHandlers は素の JSONHandler と slogctx.Handler のスループット比較。
// io.Discard に書き出して serialize 以外のオーバーヘッドを見やすくしている。
// 実行: go test -bench=. -benchmem ./src/slog/ctx/
//
// 参考値 (Apple M3 Pro / darwin/arm64 / Go 1.26):
//
//	plain                    340 ns/op     0 B/op    0 allocs/op
//	plain_group              393 ns/op     0 B/op    0 allocs/op
//	plain_logattr            370 ns/op     0 B/op    0 allocs/op  // trace_id を毎回明示的に渡す
//	slogctx                  348 ns/op     0 B/op    0 allocs/op  // wrap だけならノーコスト
//	slogctx_ctxattr          527 ns/op   256 B/op    6 allocs/op  // ctx 属性焼き付け
//	slogctx_ctxattr_group    664 ns/op   608 B/op   12 allocs/op  // + group 2 段の rebuild
//
// plain_logattr との対比で、ctx-based 注入の追加コストは +157 ns / +256 B / +6 allocs。
// 利便性と引き換えに払うコストの目安。
func BenchmarkHandlers(b *testing.B) {
	plain := slog.NewJSONHandler(io.Discard, nil)
	wrapped := NewHandler(slog.NewJSONHandler(io.Discard, nil))

	bg := context.Background()
	ctxAttr := WithAttr(bg, slog.String("trace_id", "abc123"))

	plainLogger := slog.New(plain)
	plainGroupLogger := slog.New(plain).WithGroup("g1").WithGroup("g2")
	wrappedLogger := slog.New(wrapped)
	wrappedGroupLogger := slog.New(wrapped).WithGroup("g1").WithGroup("g2")

	cases := []struct {
		name string
		do   func()
	}{
		// baseline: 素の JSONHandler
		{"plain", func() { plainLogger.InfoContext(bg, "msg", "k", "v") }},
		{"plain_group", func() { plainGroupLogger.InfoContext(bg, "msg", "k", "v") }},
		// 代替パターン: slogctx を使わず、毎回 LogAttrs で trace_id を明示的に渡す
		{"plain_logattr", func() {
			plainLogger.LogAttrs(bg, slog.LevelInfo, "msg",
				slog.String("trace_id", "abc123"),
				slog.String("k", "v"))
		}},

		// slogctx wrap のみ (ctx 属性なし) — wrap オーバーヘッド実測
		{"slogctx", func() { wrappedLogger.InfoContext(bg, "msg", "k", "v") }},

		// 主用途: ctx 属性あり (group なし / あり)
		{"slogctx_ctxattr", func() { wrappedLogger.InfoContext(ctxAttr, "msg", "k", "v") }},
		{"slogctx_ctxattr_group", func() { wrappedGroupLogger.InfoContext(ctxAttr, "msg", "k", "v") }},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				c.do()
			}
		})
	}
}

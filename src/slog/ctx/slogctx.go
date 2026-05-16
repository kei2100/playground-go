package slogctx

import (
	"context"
	"log/slog"
	"slices"
)

type contextKeyType struct{}

var contextKey = contextKeyType{}

type node struct {
	attr   slog.Attr
	parent *node
}

// WithAttr は attr を持つ新しい context を返します。
// 同じ context に対し複数回呼ぶと属性が積み重なり、Handler 経由のログ出力時に
// すべてトップレベル属性として書き出されます。親 context は変更しません。
func WithAttr(parent context.Context, attr slog.Attr) context.Context {
	n := &node{attr: attr}
	if pn, ok := parent.Value(contextKey).(*node); ok {
		n.parent = pn
	}
	return context.WithValue(parent, contextKey, n)
}

func attrsFromContext(ctx context.Context) []slog.Attr {
	n, ok := ctx.Value(contextKey).(*node)
	if !ok {
		return nil
	}

	var attrs []slog.Attr
	for {
		attrs = append(attrs, n.attr)
		if n.parent == nil {
			break
		}
		n = n.parent
	}
	slices.Reverse(attrs)
	return attrs
}

var _ slog.Handler = &Handler{}

// Handler は WithAttr で context に積まれた属性を、ラップした slog.Handler の
// トップレベル属性として出力する slog.Handler 実装です。WithGroup を任意の深さで
// 重ねても、ctx 由来の属性は group の外側 (ルート) に留まります。
//
// 実装ノート: WithGroup と group 配下の WithAttrs を遅延適用にし、Handle 時に
// 「ctx 属性をラップ先に焼く → 保留した group / 内側 attrs を順に再生する」
// という順で組み直すことで上記を実現しています。
type Handler struct {
	wrap   slog.Handler
	groups []groupState // WithGroup と「group 内 WithAttrs」を遅延保持
}

type groupState struct {
	name  string
	attrs []slog.Attr
}

// NewHandler は Handler を作成して返します。
func NewHandler(wrap slog.Handler) *Handler {
	return &Handler{wrap: wrap}
}

// Enabled はラップ先 Handler に委譲します。
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.wrap.Enabled(ctx, level)
}

// Handle は WithAttr で ctx に積まれた属性をトップレベル属性として書き出した上で、
// ラップ先 Handler にレコードの出力を委譲します。
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	wrap := h.wrap

	// ctx 由来の属性をトップレベル属性として wrap に積む
	if attrs := attrsFromContext(ctx); len(attrs) > 0 {
		wrap = wrap.WithAttrs(attrs)
	}
	// 蓄積していた group / group 内 attrs を再生
	for _, g := range h.groups {
		wrap = wrap.WithGroup(g.name)
		if len(g.attrs) > 0 {
			wrap = wrap.WithAttrs(g.attrs)
		}
	}
	return wrap.Handle(ctx, record)
}

// WithAttrs は与えられた属性を含む新しい Handler を返します。
// 元の Handler は変更しません。
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	// group がまだ無ければ wrap に即時積む
	if len(h.groups) == 0 {
		return &Handler{
			wrap:   h.wrap.WithAttrs(attrs),
			groups: h.groups,
		}
	}
	// groups があれば最後の group の attrs に積む。
	// slices.Clone は派生 Handler 間で backing array を共有しないため。
	newGroups := slices.Clone(h.groups)
	last := newGroups[len(newGroups)-1]
	last.attrs = append(slices.Clone(last.attrs), attrs...)
	newGroups[len(newGroups)-1] = last
	return &Handler{wrap: h.wrap, groups: newGroups}
}

// WithGroup は以後の属性が name 配下にネストされる新しい Handler を返します。
// WithAttr で ctx に積まれた属性はこの group の影響を受けず、ルートレベルに留まります。
// 元の Handler は変更しません。
func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &Handler{
		wrap:   h.wrap,
		groups: append(slices.Clone(h.groups), groupState{name: name}),
	}
}

package context

import (
	"context"
	"testing"
)

const someTokenKey = "__tokenKey__"

// WithSomeToken returns a copy of ctx in witch the token is set
func WithSomeToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, someTokenKey, token)
}

// GetSomeToken from ctx
func GetSomeToken(ctx context.Context) (string, bool) {
	v := ctx.Value(someTokenKey)
	token, ok := v.(string)
	return token, ok
}

func TestSomeToken(t *testing.T) {
	ctx := context.Background()

	_, ok := GetSomeToken(ctx)
	if g, w := ok, false; g != w {
		t.Errorf("ok got %v, want %v", g, w)
	}

	ctx = WithSomeToken(ctx, "token")

	token, ok := GetSomeToken(ctx)
	if g, w := ok, true; g != w {
		t.Errorf("ok got %v, want %v", g, w)
	}
	if g, w := token, "token"; g != w {
		t.Errorf("token got %v, want %v", g, w)
	}
}

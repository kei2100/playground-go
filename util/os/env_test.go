package os

import (
	"os"
	"testing"
)

func TestGetenvOrDefault(t *testing.T) {
	k := "__TEST__"
	bak := os.Getenv(k)
	defer os.Setenv(k, bak)
	os.Unsetenv(bak)

	if g, w := GetenvOrDefault(k, "default"), "default"; g != w {
		t.Errorf("value got %v, want %v", g, w)
	}

	os.Setenv(k, "set")
	if g, w := GetenvOrDefault(k, "default"), "set"; g != w {
		t.Errorf("value got %v, want %v", g, w)
	}
}

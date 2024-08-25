package go1_23

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyFS(t *testing.T) {
	t.Cleanup(func() {
		_ = os.RemoveAll("testdata/dst")
	})

	err := os.CopyFS("testdata/dst", os.DirFS("testdata/src"))
	if !assert.NoError(t, err) {
		return
	}

	if assert.DirExists(t, "testdata/dst") {
		assert.FileExists(t, "testdata/dst/empty")
	}
}

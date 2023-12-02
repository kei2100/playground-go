package ast

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectCallErrorsIs(t *testing.T) {
	detectPositions, err := DetectCallErrorsIs("testdata")
	assert.NoError(t, err)
	wantFiles, err := filepath.Glob(filepath.Join("testdata", "detect*.go"))
	assert.NoError(t, err)
	assert.Len(t, detectPositions, len(wantFiles))
	for _, pos := range detectPositions {
		fmt.Println(pos.String())
		filename := filepath.Base(pos.Filename)
		assert.Truef(t, strings.HasPrefix(filename, "detect"), "unexpected detect %s", filename)
	}
}

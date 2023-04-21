package go1_20

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultipleErrors(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	wrap1 := errors.Join(err1, err2)
	wrap2 := fmt.Errorf("foo: bar: %w, %w", err1, err2)

	assert.True(t, errors.Is(wrap1, err1))
	assert.True(t, errors.Is(wrap1, err2))
	assert.True(t, errors.Is(wrap2, err1))
	assert.True(t, errors.Is(wrap2, err2))

	wrap3 := errors.Join(nil, nil)
	assert.NoError(t, wrap3)

	wrap4 := errors.Join(err1, nil)
	assert.Error(t, wrap4)
	assert.True(t, errors.Is(wrap4, err1))

	wrap5 := errors.Join(nil, err1)
	assert.Error(t, wrap5)
	assert.True(t, errors.Is(wrap5, err1))
}

func TestHandleCloseError(t *testing.T) {
	writeFunc := func() (err error) {
		dir := t.TempDir()
		f, err := os.OpenFile(filepath.Join(dir, "tmp"), os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		defer func() {
			if err2 := f.Close(); err2 != nil {
				// Close で発生したエラーを Join して返却
				// * 大元の err が nil でなくても err2 で上書きしてしまう心配がない
				// * 大元の err が nil でも err2 の内容をエラーとして返却できる
				err = errors.Join(err, err2)
			}
		}()
		if _, err = io.WriteString(f, "hello"); err != nil {
			return err
		}
		return nil
	}
	if err := writeFunc(); err != nil {
		t.Error(err)
	}
}

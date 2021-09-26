package filepath

import (
	"path/filepath"
	"runtime"
)

func CurDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}

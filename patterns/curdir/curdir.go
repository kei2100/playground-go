package curdir

import (
	"fmt"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func printMyPath() {
	fmt.Println(basepath)
}

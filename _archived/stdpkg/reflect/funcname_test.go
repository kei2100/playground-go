package reflect

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func myFunc() {
}

func TestReflectFuncName(t *testing.T) {
	fmt.Println(runtime.FuncForPC(reflect.ValueOf(myFunc).Pointer()).Name())
}

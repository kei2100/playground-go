package concurrent

import (
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"
)

func assertABC(t *testing.T, got []string) {
	t.Helper()

	want := []string{"a", "b", "c"}
	sort.Strings(got)

	if !reflect.DeepEqual(got, want) {
		t.Errorf(" got %v, want %v", got, want)
	}
}

func TestCommonMistakes(t *testing.T) {
	// ループ変数のvはループ毎に同じインスタンスが使用される。
	// go func() { .. }して、実際にgoroutineがfmt.Println(v)するタイミングでは、vの値は変更されている可能性があり、
	// 例えば下記テストでは、gotが「c, c, c」になるケースが多い

	wg := new(sync.WaitGroup)
	vals := []string{"a", "b", "c"}
	var got []string

	for _, v := range vals {
		wg.Add(1)
		go func() {
			got = append(got, v)
			wg.Done()
		}()
	}

	wg.Wait()
	assertABC(t, got)

	// 以下のようにちょっとwaitが入っていたりすると「a, b, c」となったりする
	got = nil

	for _, v := range vals {
		wg.Add(1)
		go func() {
			got = append(got, v)
			wg.Done()
		}()
		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()
}

func TestHowToBindLoopVariable(t *testing.T) {
	wg := new(sync.WaitGroup)
	vals := []string{"a", "b", "c"}
	var got []string

	// 関数呼び出し時にbindすれば、意図通り「a, b, c」になる
	for _, v := range vals {
		wg.Add(1)
		go func(v string) {
			got = append(got, v)
			wg.Done()
		}(v)
	}

	wg.Wait()
	assertABC(t, got)

	// またはループごとに新たな変数を作り出してあげれば、意図通り「a, b, c」になる
	got = nil

	for _, v := range vals {
		wg.Add(1)
		v := v
		go func() {
			got = append(got, v)
			wg.Done()
		}()
	}

	wg.Wait()
	assertABC(t, got)
}

package concurrency

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

func TestLoopVariableAddr(t *testing.T) {
	vals := []string{"a", "b", "c", "d", "e", "f", "g"}

	// ループ変数のvはループ毎に同じインスタンスが使用される
	var sameptrs []uintptr
	for v := range vals {
		sameptrs = append(sameptrs, reflect.ValueOf(&v).Pointer())
	}
	for i := range sameptrs {
		for j := range sameptrs {
			if i == j {
				continue
			}
			if g, w := sameptrs[i], sameptrs[j]; g != w {
				t.Errorf(" got[%v] %v, want %v", i, g, w)
			}
		}
	}

	// 新たにvを宣言すれば異なるインスタンスになる
	var differentptrs []uintptr
	for v := range vals {
		v := v
		differentptrs = append(differentptrs, reflect.ValueOf(&v).Pointer())
	}
	for i := range differentptrs {
		for j := range differentptrs {
			if i == j {
				continue
			}
			if g, w := differentptrs[i], differentptrs[j]; g == w {
				t.Errorf(" got[%v] %v, want not %v", i, g, w)
			}
		}
	}
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

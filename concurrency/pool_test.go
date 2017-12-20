package concurrency

import (
	"fmt"
	"sync"
	"testing"
)

// sync.Poolは一時的なオブジェクトのPool
// groutineセーフにオブジェクトを共有するPoolからGetしたり、Putしたりすることができる。
// 内部は弱参照になっており、sync.Pool内でしか参照されていないオブジェクトは、
// 予告なしに解放される可能性がある。

type mark struct {
	used bool
}

func TestPool(t *testing.T) {
	p := &sync.Pool{
		New: func() interface{} {
			return mark{}
		},
	}
	wg := new(sync.WaitGroup)
	const n = 1000
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(num int) {
			defer wg.Done()
			b := p.Get().(mark)
			fmt.Printf("num:%v used:%v\n", num, b.used)
			b.used = true
			p.Put(b)
		}(i)
	}

	wg.Wait()
}

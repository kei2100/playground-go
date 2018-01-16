package goodparts

import (
	"fmt"
	"sync"
)

// 構造体のメンバはゼロ値に意味があるように定義するのが良い。
// そのうえで、new(struct)したら、どのようなメソッドでも動くのが良い構造体のメンバ設計、な気がしている

type GoodCache struct {
	mu    sync.Mutex
	value string
}

func (c *GoodCache) Store(v string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = v
}

func (c *GoodCache) Get() string {
	return c.value
}

type BadCache struct {
	mu    *sync.Mutex
	value string
}

func (c *BadCache) Store(v string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = v
}

func (c *BadCache) Get() string {
	return c.value
}

// ネストも可能

type LoggingGoodCache struct {
	cache GoodCache
}

func (c *LoggingGoodCache) Store(v string) {
	fmt.Printf("store %v\n", v)
	c.cache.Store(v)
}

func ExampleGoodBadCache() {
	g := new(GoodCache)
	g.Store("good")
	fmt.Println(g.Get())

	b := &BadCache{mu: new(sync.Mutex)}
	b.Store("bad")
	fmt.Println(b.Get())

	lg := new(LoggingGoodCache)
	lg.Store("good")

	// Output:
	// good
	// bad
	// store good
}

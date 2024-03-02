package rand

import (
	"errors"
	"math/rand/v2"
)

// WeightedRandSelect は重み付き乱択を行います
type WeightedRandSelect[T comparable] struct {
	sum   float64
	items map[T]float64
}

// NewWeightedRandSelect は WeightedRandSelect を生成して返却します。
// 乱択対象の items は、key にアイテム固有のキー、value にアイテムの重みを指定してください。
// アイテムの重み合計が0以下の場合、エラーを返却します。
func NewWeightedRandSelect[T comparable](items map[T]float64) (*WeightedRandSelect[T], error) {
	var sum float64
	filtered := make(map[T]float64, len(items))
	for k, w := range items {
		if w <= 0 {
			continue
		}
		sum += w
		filtered[k] = w
	}
	if sum <= 0 {
		return nil, errors.New("rand: sum of items weights must be greater than zero")
	}
	return &WeightedRandSelect[T]{
		sum:   sum,
		items: filtered,
	}, nil
}

// Select は ws に渡された重み付きアイテムの乱択を行います。
func (ws *WeightedRandSelect[T]) Select() T {
	rv := ws.randValue()
	for k, w := range ws.items {
		rv = rv - w
		if rv <= 0 {
			return k
		}
	}
	panic("unreachable")
}

func (ws *WeightedRandSelect[T]) randValue() float64 {
	// https://stackoverflow.com/questions/49746992/generate-random-float64-numbers-in-specific-range-using-golang
	return rand.Float64() * ws.sum
}

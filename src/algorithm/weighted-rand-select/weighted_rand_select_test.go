package rand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightedRandSelect_Select(t *testing.T) {
	items := map[string]float64{
		"red":   0,
		"green": 0.70,
		"blue":  0.30,
	}
	ws, err := NewWeightedRandSelect(items)
	if err != nil {
		t.Fatal(err)
	}

	const n = 10000
	var red, green, blue int
	for i := 0; i < n; i++ {
		_ = i
		got := ws.Select()
		switch got {
		case "red":
			red++
		case "green":
			green++
		case "blue":
			blue++
		}
	}
	const deltaRate = 0.03
	assert.Equal(t, 0, red)
	assert.InDelta(t, n*items["green"], green, n*deltaRate)
	assert.InDelta(t, n*items["blue"], blue, n*deltaRate)
}

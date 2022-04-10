package nradix

import (
	"math"
)

// NRadix converts a number to specified radix string
type NRadix struct {
	letters []rune
}

// New creates a NRadix
func New(radixLetters string) *NRadix {
	letters := []rune(radixLetters)
	if len(letters) < 2 {
		panic("radixLetters length must greater than 2")
	}
	return &NRadix{
		letters: letters,
	}
}

// ConvertToString converts an integer number to specified radix string
func (nr *NRadix) ConvertToString(num int64) string {
	var negative bool
	if num < 0 {
		negative = true
		num = int64(math.Abs(float64(num)))
	}
	n := int64(len(nr.letters))
	l := calcStringLen(num, n)
	ret := make([]rune, l)
	for i := l - 1; i >= 0; i-- {
		q, r := num/n, num%n
		ret[i] = nr.letters[r]
		num = q
	}
	if negative {
		return "-" + string(ret)
	}
	return string(ret)
}

func calcStringLen(num, n int64) int64 {
	var l int64
	for {
		num = num / n
		if num == 0 {
			break
		}
		l++
	}
	return l + 1
}

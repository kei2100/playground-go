package maptest

import (
	"log"
	"testing"
)

func TestNilMap(t *testing.T) {
	var m map[int]int
	// m[0] = 0 // panics

	m = make(map[int]int)
	m[0] = 0
}

func TestRef(t *testing.T) {
	m := map[int]int{0: 0}
	mm := m
	mm[1] = 1

	log.Println(m[0])
	log.Println(m[1])
}

package fmt

import (
	"fmt"
	"log"
	"testing"
)

func TestPadding(t *testing.T) {
	lpad := fmt.Sprintf("% 10s", "foo")
	if g, w := lpad, "       foo"; g != w {
		t.Errorf("\ngot:\n%v\nwant:\n%v", g, w)
	}

	rpad := fmt.Sprintf("%- 10s", "foo")
	if g, w := rpad, "foo       "; g != w {
		t.Errorf("\ngot:\n%v\nwant:\n%v", g, w)
	}
}

type custom struct {
}

func (custom) Format(s fmt.State, c rune) {
	log.Println(s.Width())
	log.Println(s.Precision())
	log.Println(s.Flag('-'))
	log.Println(string(c))
}

func TestCustom(t *testing.T) {
	c := custom{}
	fmt.Printf("%-010.2s", c)
	// Width: 10
	// Precision: 2
	// s.Flag('-'): true
	// runeは「's'」
	//
	// %-0の「0」の情報を取得する方法がよくわからない...
}

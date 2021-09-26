package norm

import (
	"log"
	"testing"

	"golang.org/x/text/unicode/norm"
)

func TestNormalize(t *testing.T) {
	s := "ｶﾞ"

	nfkced := norm.NFKC.String(s)
	nfkded := norm.NFKD.String(s)

	log.Printf("NFKCed:%v NFKDed:%v", nfkced, nfkded)

	if nfkced == nfkded {
		t.Error("expect not equal")
	}
	if norm.NFD.String(nfkced) != nfkded {
		t.Error("expect equal")
	}
}

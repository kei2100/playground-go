package rune

import (
	"log"

	"golang.org/x/text/width"
)

// IsHankaku reports whether the rune is half-width or not
func IsHankaku(r rune) bool {
	k := width.LookupRune(r).Kind()
	switch k {
	case width.Neutral: // 東アジアの組版には通常出現せず、全角でも半角でもない。アラビア文字など。
		return false
	case width.EastAsianAmbiguous: // 文脈によって文字幅が異なる文字。東アジアの組版とそれ以外の組版の両方に出現し、東アジアの従来文字コードではいわゆる全角として扱われることがある。ギリシア文字やキリル文字など。
		return false
	case width.EastAsianWide: // EastAsianFullwidth ではない全角文字。漢字など。
		return false
	case width.EastAsianNarrow: // EastAsianHalfwidth ではない、対応する全角文字が存在する半角文字。半角英数字など。
		return true
	case width.EastAsianFullwidth: // 互換分解特性を持つ全角文字。
		return false
	case width.EastAsianHalfwidth: // 互換分解特性を持つ半角文字。半角カナなど。
		return true
	default:
		log.Printf("[WARN] unexpected kind of the text width %s", k.String())
		return false
	}
}

// IsZenkaku reports whether the rune is full-width or not
func IsZenkaku(r rune) bool {
	k := width.LookupRune(r).Kind()
	switch k {
	case width.Neutral: // 東アジアの組版には通常出現せず、全角でも半角でもない。アラビア文字など。
		return false
	case width.EastAsianAmbiguous: // 文脈によって文字幅が異なる文字。東アジアの組版とそれ以外の組版の両方に出現し、東アジアの従来文字コードではいわゆる全角として扱われることがある。ギリシア文字やキリル文字など。
		return false
	case width.EastAsianWide: // EastAsianFullwidth ではない全角文字。漢字など。
		return true
	case width.EastAsianNarrow: // EastAsianHalfwidth ではない、対応する全角文字が存在する半角文字。半角英数字など。
		return false
	case width.EastAsianFullwidth: // 互換分解特性を持つ全角文字。
		return true
	case width.EastAsianHalfwidth: // 互換分解特性を持つ半角文字。半角カナなど。
		return false
	default:
		log.Printf("[WARN] unexpected kind of the text width %s", k.String())
		return false
	}
}

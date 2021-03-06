package width

import (
	"strings"
	"testing"

	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
)

func toHalf(s string) string {
	// 全角カタカナの濁点半濁点付きは合成列（Composite Sequence）として変換する。
	// NFDする代わり
	rep := strings.NewReplacer(
		"ガ", "ｶﾞ",
		"ギ", "ｷﾞ",
		"グ", "ｸﾞ",
		"ゲ", "ｹﾞ",
		"ゴ", "ｺﾞ",
		"ザ", "ｻﾞ",
		"ジ", "ｼﾞ",
		"ズ", "ｽﾞ",
		"ゼ", "ｾﾞ",
		"ゾ", "ｿﾞ",
		"ダ", "ﾀﾞ",
		"ヂ", "ﾁﾞ",
		"ヅ", "ﾂﾞ",
		"デ", "ﾃﾞ",
		"ド", "ﾄﾞ",
		"バ", "ﾊﾞ",
		"パ", "ﾊﾟ",
		"ビ", "ﾋﾞ",
		"ピ", "ﾋﾟ",
		"ブ", "ﾌﾞ",
		"プ", "ﾌﾟ",
		"ベ", "ﾍﾞ",
		"ペ", "ﾍﾟ",
		"ボ", "ﾎﾞ",
		"ポ", "ﾎﾟ",
		"ヴ", "ｳﾞ",
		"ヷ", "ﾜﾞ",
		"ヸ", "ヰﾞ", // 非濁点部半角対応なし
		"ヹ", "ヱﾞ", // 非濁点部半角対応なし
		"ヺ", "ｦﾞ",
		"゛", "ﾞ",
		"゜", "ﾟ",
	)
	s = rep.Replace(s)
	return width.Narrow.String(s)
}

func toFull(s string) string {
	// 半角カタカナの濁点半濁点付きは事前合成系（Precomposed Character）として変換する。
	// NFKCする代わり
	rep := strings.NewReplacer(
		"ｶﾞ", "ガ",
		"ｷﾞ", "ギ",
		"ｸﾞ", "グ",
		"ｹﾞ", "ゲ",
		"ｺﾞ", "ゴ",
		"ｻﾞ", "ザ",
		"ｼﾞ", "ジ",
		"ｽﾞ", "ズ",
		"ｾﾞ", "ゼ",
		"ｿﾞ", "ゾ",
		"ﾀﾞ", "ダ",
		"ﾁﾞ", "ヂ",
		"ﾂﾞ", "ヅ",
		"ﾃﾞ", "デ",
		"ﾄﾞ", "ド",
		"ﾊﾞ", "バ",
		"ﾊﾟ", "パ",
		"ﾋﾞ", "ビ",
		"ﾋﾟ", "ピ",
		"ﾌﾞ", "ブ",
		"ﾌﾟ", "プ",
		"ﾍﾞ", "ベ",
		"ﾍﾟ", "ペ",
		"ﾎﾞ", "ボ",
		"ﾎﾟ", "ポ",
		"ｳﾞ", "ヴ",
		"ﾜﾞ", "ヷ",
		"ヰﾞ", "ヸ", // 非濁点部半角対応なし
		"ヱﾞ", "ヹ", // 非濁点部半角対応なし
		"ｦﾞ", "ヺ",
		"ﾞ", "゛",
		"ﾟ", "゜",
	)
	s = rep.Replace(s)
	return width.Widen.String(s)
}

func TestToHalf(t *testing.T) {
	tt := [][2]string{
		{"アイウ", "ｱｲｳ"},
		{norm.NFC.String("ガ"), "ｶﾞ"},
		{norm.NFD.String("ガ"), "ｶﾞ"},
		{"あいう", "あいう"},
		{"神と神", "神と神"}, // 神がnormalizeされないこと
		{"あいうアイウえオ", "あいうｱｲｳえｵ"},
		{"ㇰ", "ㇰ"}, // Unicodeの片仮名拡張.対応する半角文字なし
		{"「」", "｢｣"},
		{"、。", "､｡"},
		// カタカナ一覧
		{"ァ", "ｧ"},
		{"ア", "ｱ"},
		{"ィ", "ｨ"},
		{"イ", "ｲ"},
		{"ゥ", "ｩ"},
		{"ウ", "ｳ"},
		{"ェ", "ｪ"},
		{"エ", "ｴ"},
		{"ォ", "ｫ"},
		{"オ", "ｵ"},
		{"カ", "ｶ"},
		{"ガ", "ｶﾞ"},
		{"キ", "ｷ"},
		{"ギ", "ｷﾞ"},
		{"ク", "ｸ"},
		{"グ", "ｸﾞ"},
		{"ケ", "ｹ"},
		{"ゲ", "ｹﾞ"},
		{"コ", "ｺ"},
		{"ゴ", "ｺﾞ"},
		{"サ", "ｻ"},
		{"ザ", "ｻﾞ"},
		{"シ", "ｼ"},
		{"ジ", "ｼﾞ"},
		{"ス", "ｽ"},
		{"ズ", "ｽﾞ"},
		{"セ", "ｾ"},
		{"ゼ", "ｾﾞ"},
		{"ソ", "ｿ"},
		{"ゾ", "ｿﾞ"},
		{"タ", "ﾀ"},
		{"ダ", "ﾀﾞ"},
		{"チ", "ﾁ"},
		{"ヂ", "ﾁﾞ"},
		{"ッ", "ｯ"},
		{"ツ", "ﾂ"},
		{"ヅ", "ﾂﾞ"},
		{"テ", "ﾃ"},
		{"デ", "ﾃﾞ"},
		{"ト", "ﾄ"},
		{"ド", "ﾄﾞ"},
		{"ナ", "ﾅ"},
		{"ニ", "ﾆ"},
		{"ヌ", "ﾇ"},
		{"ネ", "ﾈ"},
		{"ノ", "ﾉ"},
		{"ハ", "ﾊ"},
		{"バ", "ﾊﾞ"},
		{"パ", "ﾊﾟ"},
		{"ヒ", "ﾋ"},
		{"ビ", "ﾋﾞ"},
		{"ピ", "ﾋﾟ"},
		{"フ", "ﾌ"},
		{"ブ", "ﾌﾞ"},
		{"プ", "ﾌﾟ"},
		{"ヘ", "ﾍ"},
		{"ベ", "ﾍﾞ"},
		{"ペ", "ﾍﾟ"},
		{"ホ", "ﾎ"},
		{"ボ", "ﾎﾞ"},
		{"ポ", "ﾎﾟ"},
		{"マ", "ﾏ"},
		{"ミ", "ﾐ"},
		{"ム", "ﾑ"},
		{"メ", "ﾒ"},
		{"モ", "ﾓ"},
		{"ャ", "ｬ"},
		{"ヤ", "ﾔ"},
		{"ュ", "ｭ"},
		{"ユ", "ﾕ"},
		{"ョ", "ｮ"},
		{"ヨ", "ﾖ"},
		{"ラ", "ﾗ"},
		{"リ", "ﾘ"},
		{"ル", "ﾙ"},
		{"レ", "ﾚ"},
		{"ロ", "ﾛ"},
		{"ヮ", "ヮ"}, // 半角対応なし
		{"ワ", "ﾜ"},
		{"ヰ", "ヰ"}, // 半角対応なし
		{"ヱ", "ヱ"}, // 半角対応なし
		{"ヲ", "ｦ"},
		{"ン", "ﾝ"},
		{"ヴ", "ｳﾞ"},
		{"ヵ", "ヵ"}, // 半角対応なし
		{"ヶ", "ヶ"}, // 半角対応なし
		{"ヷ", "ﾜﾞ"},
		{"ヸ", "ヰﾞ"}, // 非濁点部半角対応なし
		{"ヹ", "ヱﾞ"}, // 非濁点部半角対応なし
		{"ヺ", "ｦﾞ"},
		{"゛", "ﾞ"},
		{"゜", "ﾟ"},
		{"", ""},
	}
	for _, te := range tt {
		if g, w := toHalf(te[0]), te[1]; g != w {
			t.Errorf(" got %v, want %v", g, w)
		}
	}
}

func TestToFull(t *testing.T) {
	tt := [][2]string{
		{"ｱｲｳ", "アイウ"},
		{"ｶﾞ", norm.NFC.String("ガ")},
		{"あいう", "あいう"},
		{"神と神", "神と神"}, // 神がnormalizeされないこと
		{"あいうｱｲｳえｵ", "あいうアイウえオ"},
		{"ㇰ", "ㇰ"}, // Unicodeの片仮名拡張.対応する半角文字なし
		{"｢｣", "「」"},
		{"､｡", "、。"},
		// カタカナ一覧
		{"ｱ", "ア"},
		{"ｨ", "ィ"},
		{"ｲ", "イ"},
		{"ｩ", "ゥ"},
		{"ｳ", "ウ"},
		{"ｪ", "ェ"},
		{"ｴ", "エ"},
		{"ｫ", "ォ"},
		{"ｵ", "オ"},
		{"ｶ", "カ"},
		{"ｶﾞ", "ガ"},
		{"ｷ", "キ"},
		{"ｷﾞ", "ギ"},
		{"ｸ", "ク"},
		{"ｸﾞ", "グ"},
		{"ｹ", "ケ"},
		{"ｹﾞ", "ゲ"},
		{"ｺ", "コ"},
		{"ｺﾞ", "ゴ"},
		{"ｻ", "サ"},
		{"ｻﾞ", "ザ"},
		{"ｼ", "シ"},
		{"ｼﾞ", "ジ"},
		{"ｽ", "ス"},
		{"ｽﾞ", "ズ"},
		{"ｾ", "セ"},
		{"ｾﾞ", "ゼ"},
		{"ｿ", "ソ"},
		{"ｿﾞ", "ゾ"},
		{"ﾀ", "タ"},
		{"ﾀﾞ", "ダ"},
		{"ﾁ", "チ"},
		{"ﾁﾞ", "ヂ"},
		{"ｯ", "ッ"},
		{"ﾂ", "ツ"},
		{"ﾂﾞ", "ヅ"},
		{"ﾃ", "テ"},
		{"ﾃﾞ", "デ"},
		{"ﾄ", "ト"},
		{"ﾄﾞ", "ド"},
		{"ﾅ", "ナ"},
		{"ﾆ", "ニ"},
		{"ﾇ", "ヌ"},
		{"ﾈ", "ネ"},
		{"ﾉ", "ノ"},
		{"ﾊ", "ハ"},
		{"ﾊﾞ", "バ"},
		{"ﾊﾟ", "パ"},
		{"ﾋ", "ヒ"},
		{"ﾋﾞ", "ビ"},
		{"ﾋﾟ", "ピ"},
		{"ﾌ", "フ"},
		{"ﾌﾞ", "ブ"},
		{"ﾌﾟ", "プ"},
		{"ﾍ", "ヘ"},
		{"ﾍﾞ", "ベ"},
		{"ﾍﾟ", "ペ"},
		{"ﾎ", "ホ"},
		{"ﾎﾞ", "ボ"},
		{"ﾎﾟ", "ポ"},
		{"ﾏ", "マ"},
		{"ﾐ", "ミ"},
		{"ﾑ", "ム"},
		{"ﾒ", "メ"},
		{"ﾓ", "モ"},
		{"ｬ", "ャ"},
		{"ﾔ", "ヤ"},
		{"ｭ", "ュ"},
		{"ﾕ", "ユ"},
		{"ｮ", "ョ"},
		{"ﾖ", "ヨ"},
		{"ﾗ", "ラ"},
		{"ﾘ", "リ"},
		{"ﾙ", "ル"},
		{"ﾚ", "レ"},
		{"ﾛ", "ロ"},
		{"ヮ", "ヮ"}, // 半角対応なし
		{"ﾜ", "ワ"},
		{"ヰ", "ヰ"}, // 半角対応なし
		{"ヱ", "ヱ"}, // 半角対応なし
		{"ｦ", "ヲ"},
		{"ﾝ", "ン"},
		{"ｳﾞ", "ヴ"},
		{"ヵ", "ヵ"}, // 半角対応なし
		{"ヶ", "ヶ"}, // 半角対応なし
		{"ﾜﾞ", "ヷ"},
		{"ヰﾞ", "ヸ"}, // 非濁点部半角対応なし
		{"ヱﾞ", "ヹ"}, // 非濁点部半角対応なし
		{"ｦﾞ", "ヺ"},
		{"ﾞ", "゛"},
		{"ﾟ", "゜"},
		{"", ""},
	}
	for _, te := range tt {
		if g, w := toFull(te[0]), te[1]; g != w {
			t.Errorf(" got %v, want %v", g, w)
		}
	}
}

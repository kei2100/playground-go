package mail_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kei2100/playground-go/src/mail"
)

func TestIsRFC5321EmailAddress(t *testing.T) {
	tt := []struct {
		email string
		valid bool
	}{
		// --- mailBox ---
		{"foo@example.com", true},
		{" foo@example.com", false},
		{"foo @example.com", false},
		{"foo@ example.com", false},
		{"foo@example.com ", false},
		{"foo@@example.com", false},
		{"<foo@example.com>", false},
		{"foo <foo@example.com>", false},
		// --- local-part ---
		// local-part / dot-string
		// allowed chars
		{"azAZ09!#$%&'*+-/=?^_`{|}~@example.com", true},
		{"azAZ09!#$%&'*+-/=?^_`{|}~,@example.com", false},
		{"ａzAZ09!#$%&'*+-/=?^_`{|}~@example.com", false}, // Full-width ａ
		// dots
		{"a@example.com", true},
		{"a.a@example.com", true},
		{"a.a.a@example.com", true},
		{"a.a.@example.com", false},
		{".a.a@example.com", false},
		{"a..a@example.com", false},
		// local-part / Quoted-string
		// qtextSMTP
		{`"` + stringsByCodePointsRange(32, 33) + `"@example.com`, true},
		{`"` + stringsByCodePointsRange(35, 91) + `"@example.com`, true},
		{`"` + stringsByCodePointsRange(93, 126) + `"@example.com`, true},
		{`"` + stringsByCodePointsRange(31, 33) + `"@example.com`, false},
		{`"` + stringsByCodePointsRange(32, 34) + `"@example.com`, false},
		{`"` + stringsByCodePointsRange(35, 92) + `"@example.com`, false},
		{`"` + stringsByCodePointsRange(93, 127) + `"@example.com`, false},
		// quoted-pairSMTP
		{`"` + escape(stringsByCodePointsRange(32, 34), '\\') + `"@example.com`, true},
		// %d35-92
		{`"` + escape(stringsByCodePointsRange(35, 65), '\\') + `"@example.com`, true},
		{`"` + escape(stringsByCodePointsRange(65, 92), '\\') + `"@example.com`, true},
		// %d93-126
		{`"` + escape(stringsByCodePointsRange(93, 123), '\\') + `"@example.com`, true},
		{`"` + escape(stringsByCodePointsRange(124, 126), '\\') + `"@example.com`, true},
		// out of range
		{`"` + escape(stringsByCodePointsRange(31, 34), '\\') + `"@example.com`, false},
		{`"` + escape(stringsByCodePointsRange(124, 127), '\\') + `"@example.com`, false},
		{`"\""@example.com`, true},
		{`"\\""@example.com`, false},
		// QcontentSMTP = qtextSMTP / quoted-pairSMTP
		{`" \",\\~"@example.com`, true},
		// DQUOTE
		{`" "@example.com`, true},
		{` " "@example.com`, false},
		{`" " @example.com`, false},
		// --- Domain  ---
		{"foo@a.b", true},
		{"foo@a", false}, // prevent Dotless Domain Name
		{"foo@1.1", true},
		{"foo@aa.bb", true},
		{"foo@aa. b", false},
		// sub-domain = Let-dig [Ldh-str]
		{"foo@a-a.b", true},
		{"foo@ab-cd.ef-gh", true},
		{"foo@a--a.b", true},
		{"foo@1-2-3.a-b-c", true},
		{"foo@-2-3.a-b-c", false},
		{"foo@1-2-3.a-b-", false},
		// dots
		{"foo@.a.b", false},
		{"foo@a.b.", false},
		{"foo@a..b", false},
		// --- maximum octet ---
		// local-part
		{strings.Repeat("a", 64) + "@example.com", true},
		{strings.Repeat("a", 65) + "@example.com", false},
		{`"` + strings.Repeat("@", 62) + `"@example.com`, true},
		{`"` + strings.Repeat("@", 63) + `"@example.com`, false},
		// Domain
		{"a@a." + strings.Repeat("a", 253), true},
		{"a@a." + strings.Repeat("a", 254), false},
	}
	for i, te := range tt {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := mail.IsRFC5321EmailAddress(te.email)
			if g, w := got, te.valid; g != w {
				t.Errorf("%s\ngot :%v\nwant:%v", te.email, got, te.valid)
			}
		})
	}
}

func stringsByCodePointsRange(decimalCodePointFrom, To int) string {
	var b strings.Builder
	for ; decimalCodePointFrom <= To; decimalCodePointFrom++ {
		b.WriteRune(rune(decimalCodePointFrom))
	}
	return b.String()
}

func escape(s string, escapeChar rune) string {
	var b strings.Builder
	for _, r := range s {
		b.WriteRune(escapeChar)
		b.WriteRune(r)
	}
	return b.String()
}

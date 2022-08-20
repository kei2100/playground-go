package mail

import (
	"fmt"
	"regexp"
	"strings"
)

// rfc5321EmailAddressRegexp represents a regexp of the [RFC 5321] Mailbox format.
// The whole regular expression is as follows:
//
//	\A(?:[a-zA-Z0-9!#\$%&'\*\+\-\/=\?\^_`\{\|\}~]{1,64}(?:\.[a-zA-Z0-9!#\$%&'\*\+\-\/=\?\^_`\{\|\}~]{1,64})*|"(?:[ -!#-\[\]-~]|\\[ -~]){1,64}")@[a-zA-Z0-9](?:[a-zA-Z0-9\-]{0,253}[a-zA-Z0-9])?(\.[a-zA-Z0-9](?:[a-zA-Z0-9\-]{0,253}[a-zA-Z0-9])?)+\z
//
// NOTE:
//   - Dotless domain names are not allowed.
//   - address-literal not supported yet (Mailbox = Local-part "@" ( Domain / address-literal ))
//
// [RFC 5321]: https://www.rfc-editor.org/rfc/rfc5321
var rfc5321EmailAddressRegexp = regexp.MustCompile(mailBox)

var (
	// RFC 5321 Mailbox ABNF:
	//
	//   Mailbox          = Local-part "@" ( Domain / address-literal )
	//
	//   Local-part       = Dot-string / Quoted-string
	//   Dot-string       = Atom *("."  Atom)
	//   Atom             = 1*atext
	//   ; atext refs RFC 5322
	//   atext            = ALPHA / DIGIT /    ; Printable US-ASCII
	//                      "!" / "#" /        ;  characters not including
	//                      "$" / "%" /        ;  specials.  Used for atoms.
	//                      "&" / "'" /
	//                      "*" / "+" /
	//                      "-" / "/" /
	//                      "=" / "?" /
	//                      "^" / "_" /
	//                      "`" / "{" /
	//                      "|" / "}" /
	//                      "~"
	//   Quoted-string    = DQUOTE *QcontentSMTP DQUOTE
	//   QcontentSMTP     = qtextSMTP / quoted-pairSMTP
	//   qtextSMTP        = %d32-33 /
	//                      %d35-91 /
	//                      %d93-126
	//                      ; i.e., within a quoted string, any
	//                      ; ASCII graphic or space is permitted
	//                      ; without backslash-quoting except
	//                      ; double-quote and the backslash itself.
	//   quoted-pairSMTP  = %d92 %d32-126
	//                      ; i.e., backslash followed by any ASCII
	//                      ; graphic (including itself) or SPace
	//
	//   Domain           = sub-domain *("." sub-domain)
	//   sub-domain       = Let-dig [Ldh-str]
	//   Let-dig          = ALPHA / DIGIT
	//   Ldh-str          = *( ALPHA / DIGIT / "-" ) Let-dig
	//

	// Mailbox
	mailBox = `\A` + localPart + `@` + domain + `\z`
	// local-part
	localPart = `(?:` + dotString + `|` + quotedString + `)`
	// dot-string
	dotString = atom + `(?:\.` + atom + `)*`
	atom      = atext + `{1,64}` // 64: local-part maximum octet
	atext     = `[` + alpha + digit + symbols + `]`
	alpha     = "a-zA-Z"
	digit     = "0-9"
	symbols   = `!` + `#` +
		`\$` + `%` +
		`&` + `'` +
		`\*` + `\+` +
		`\-` + `\/` +
		`=` + `\?` +
		`\^` + `_` +
		"`" + `\{` +
		`\|` + `\}` +
		`~`
	// quoted-string
	quotedString = `"` + qcontentSMTP + `{1,64}"`
	qcontentSMTP = `(?:` + qtextSMTP + `|` + quotedPairSMTP + `)`
	qtextSMTP    = `[` +
		` -!` + // %d32-33
		`#-\[` + // %d35-91
		`\]-~` + // %d93-126
		`]`
	quotedPairSMTP = `\\[ -~]` // %d32-126

	// Domain
	domain    = subDomain + `(\.` + subDomain + `)+` // prevent Dotless domain names
	subDomain = letDig + `(?:` + ldhStr + `)?`
	letDig    = `[` + alpha + digit + `]`
	ldhStr    = `[` + alpha + digit + `\-]{0,253}` + letDig // 253: Domain maximum octet is 255. and excludes leading/ending let-dig => 253
)

// IsRFC5321EmailAddress reports whether s represents RFC 5321 email address (a.k.a Mailbox)
func IsRFC5321EmailAddress(s string) bool {
	if !rfc5321EmailAddressRegexp.MatchString(s) {
		return false
	}
	at := strings.LastIndex(s, "@")
	if at <= 0 || at+1 >= len(s) {
		panic(fmt.Sprintf("unexpected email format %s\nregexp %s\n", s, rfc5321EmailAddressRegexp.String()))
	}
	localPart := s[:at]
	domain := s[at+1:]
	// https://www.rfc-editor.org/rfc/rfc5321#section-4.5.3.1.1
	// https://www.rfc-editor.org/rfc/rfc5321#section-4.5.3.1.2
	return len(localPart) <= 64 && len(domain) <= 255
}

package encoding

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// Encoding type
type Encoding string

// Available Encoding values
const (
	UnknownEncoding Encoding = "unknown"
	UTF8            Encoding = "UTF-8"
	EUCJP           Encoding = "EUC-JP"
	ShiftJIS        Encoding = "Shift_JIS"
	ISO2022JP       Encoding = "ISO-2022-JP"
)

// NewWriter returns a new Writer that wraps w
//
//	wに書き込まれたUTF8を、このEncodingでエンコードするWriterを返却する。
func (e Encoding) NewWriter(w io.Writer) io.WriteCloser {
	return transform.NewWriter(w, e.encoding().NewEncoder())
}

// NewReader returns a new Reader that wraps r
//
//	このEncodingをUTF8にデコードするReaderを返却する。
func (e Encoding) NewReader(r io.Reader) io.Reader {
	return transform.NewReader(r, e.encoding().NewDecoder())
}

func (e Encoding) encoding() encoding.Encoding {
	switch e {
	case UTF8:
		return encoding.Nop
	case EUCJP:
		return japanese.EUCJP
	case ShiftJIS:
		return japanese.ShiftJIS
	case ISO2022JP:
		return japanese.ISO2022JP
	default:
		return encoding.Nop
	}
}

func TestEncoding(t *testing.T) {
	var wb bytes.Buffer
	w := ShiftJIS.NewWriter(&wb)
	w.Write([]byte("あ"))

	if g, w := wb.Bytes(), []byte{130, 160}; !reflect.DeepEqual(g, w) {
		t.Errorf(" got %v, want %v", g, w)
	}

	rb := make([]byte, len([]byte("あ")))
	r := ShiftJIS.NewReader(&wb)
	r.Read(rb)
	if g, w := rb, []byte("あ"); !reflect.DeepEqual(g, w) {
		t.Errorf(" got %v, want %v", g, w)
	}
}

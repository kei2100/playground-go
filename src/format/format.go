package format

import "go/format"

func Format(src []byte) ([]byte, error) {
	formatted, err := format.Source(src)
	return formatted, err
}

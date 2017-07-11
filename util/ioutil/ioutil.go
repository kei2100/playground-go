package ioutil

import (
	"io"
	goioutil "io/ioutil"
)

// MustReadAll reads from r until EOF and returns the data it read, or panics.
func MustReadAll(r io.Reader) []byte {
	b, err := goioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return b
}

package xml

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

const raw = `
<?xml version="1.0" encoding="UTF-8"?>
<Message>
  <Level1>
    <Level2>
      <Foo>foo foo</Foo>
    </Level2>
	<Bar>bar 1</Bar>
	<Bar>bar 2</Bar>
  </Level1>
</Message>
`

type data struct {
	Foo string   `xml:"Level1>Level2>Foo"`
	Bar []string `xml:"Level1>Bar"`
}

func TestXML(t *testing.T) {
	dec := xml.NewDecoder(bytes.NewBufferString(raw))
	var d data
	if err := dec.Decode(&d); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "foo foo", d.Foo)
	if assert.Len(t, d.Bar, 2) {
		assert.Equal(t, "bar 1", d.Bar[0])
		assert.Equal(t, "bar 2", d.Bar[1])
	}
}

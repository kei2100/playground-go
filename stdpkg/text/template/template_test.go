package template

import (
	"bytes"
	"testing"
	"text/template"
)

// ### Overview
//
// - テンプレートへのインプットはUTF-8のテキスト
// - parseした後のtemplateはparallelに実行可能
// -

func TestTextAndSpaces(t *testing.T) {
	var tmpl *template.Template
	var b *bytes.Buffer

	// default trimming
	tmpl, _ = template.New("test").Parse("{{ 23 }} < {{ 45 }}")
	b = &bytes.Buffer{}
	tmpl.Execute(b, nil)
	if g, w := b.String(), "23 < 45"; g != w {
		t.Errorf("got '%v', want '%v'", g, w)
	}

	tmpl, _ = template.New("test").Parse("{{  23  }} < {{  45  }}")
	b = &bytes.Buffer{}
	tmpl.Execute(b, nil)
	if g, w := b.String(), "23 < 45"; g != w {
		t.Errorf("got '%v', want '%v'", g, w)
	}

	tmpl, _ = template.New("test").Parse("{{		23		}} < {{ 45 }}") // tab
	b = &bytes.Buffer{}
	tmpl.Execute(b, nil)
	if g, w := b.String(), "23 < 45"; g != w {
		t.Errorf("got '%v', want '%v'", g, w)
	}

	// right trimming
	tmpl, _ = template.New("test").Parse("{{ 23 -}} < {{ 45 }}")
	b = &bytes.Buffer{}
	tmpl.Execute(b, nil)
	if g, w := b.String(), "23< 45"; g != w {
		t.Errorf("got '%v', want '%v'", g, w)
	}

	// left and right trimming
	tmpl, _ = template.New("test").Parse("{{ 23 -}} < {{- 45 }}")
	b = &bytes.Buffer{}
	tmpl.Execute(b, nil)
	if g, w := b.String(), "23<45"; g != w {
		t.Errorf("got '%v', want '%v'", g, w)
	}

	// newline trimming
	tmpl, _ = template.New("test").Parse("{{ 23 -}}\n<\r\n{{- 45 }}")
	b = &bytes.Buffer{}
	tmpl.Execute(b, nil)
	if g, w := b.String(), "23<45"; g != w {
		t.Errorf("got '%v', want '%v'", g, w)
	}
}

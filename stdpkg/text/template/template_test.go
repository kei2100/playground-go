package template

import (
	"bytes"
	"strings"
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

func TestActions(t *testing.T) {
	var tmpl *template.Template
	var b *bytes.Buffer

	t.Run("comments", func(t *testing.T) {
		// comments
		tmpl, _ = template.New("t").Parse(" {{/* comment */}} ")
		b = &bytes.Buffer{}
		tmpl.Execute(b, nil)
		if g, w := b.String(), "  "; g != w {
			t.Errorf("got '%v', want '%v'", g, w)
		}

		// with trimming
		tmpl, _ = template.New("t").Parse(" {{- /* comment */ -}} ")
		b = &bytes.Buffer{}
		tmpl.Execute(b, nil)
		if g, w := b.String(), ""; g != w {
			t.Errorf("got '%v', want '%v'", g, w)
		}
	})

	t.Run("if", func(t *testing.T) {
		tmpl, _ = template.New("t").Parse("{{- if .Message -}} present! {{- else -}} empty! {{- end -}}") // if .Message is empty

		b = &bytes.Buffer{}
		tmpl.Execute(b, map[string]string{"Message": "foo"})
		if g, w := b.String(), "present!"; g != w {
			t.Errorf("got '%v', want '%v'", g, w)
		}

		b = &bytes.Buffer{}
		tmpl.Execute(b, map[string]string{"Message": ""})
		if g, w := b.String(), "empty!"; g != w {
			t.Errorf("got '%v', want '%v'", g, w)
		}
	})

	t.Run("range", func(t *testing.T) {
		ts := strings.TrimSpace(`
{{ range $k, $v := . }}
	{{- printf "Key:%s" $k }} {{ printf "Value:%s" $v }}
{{ else }}
	{{- "empty!" }}
{{ end }}
`)
		tmpl, _ = template.New("t").Parse(ts)
		b = &bytes.Buffer{}
		tmpl.Execute(b, map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		})

		want := `
Key:key1 Value:value1
Key:key2 Value:value2
Key:key3 Value:value3
`
		if g, w := strings.TrimSpace(b.String()), strings.TrimSpace(want); g != w {
			t.Errorf("got '%v', want '%v'", g, w)
		}
	})
}

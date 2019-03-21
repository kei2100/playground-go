package template

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"text/template"
	"time"
)

// - テンプレートへのインプットはUTF-8のテキスト
// - parseした後のtemplateはparallelに実行可能

func TestTextAndSpaces(t *testing.T) {
	// default trimming
	tmpl := mustParse(t, "{{ 23 }} < {{ 45 }}")
	result := mustExecute(t, tmpl, nil)
	if g, w := result, "23 < 45"; g != w {
		t.Errorf("got '%v', want '%v'", g, w)
	}

	tmpl = mustParse(t, "{{  23  }} < {{  45  }}")
	result = mustExecute(t, tmpl, nil)
	if g, w := result, "23 < 45"; g != w {
		t.Errorf("got '%v', want '%v'", g, w)
	}

	tmpl = mustParse(t, "{{		23		}} < {{ 45 }}") // tab
	result = mustExecute(t, tmpl, nil)
	if g, w := result, "23 < 45"; g != w {
		t.Errorf("got '%v', want '%v'", g, w)
	}

	// right trimming
	tmpl = mustParse(t, "{{ 23 -}} < {{ 45 }}")
	result = mustExecute(t, tmpl, nil)
	if g, w := result, "23< 45"; g != w {
		t.Errorf("got '%v', want '%v'", g, w)
	}

	// left and right trimming
	tmpl = mustParse(t, "{{ 23 -}} < {{- 45 }}")
	result = mustExecute(t, tmpl, nil)
	if g, w := result, "23<45"; g != w {
		t.Errorf("got '%v', want '%v'", g, w)
	}

	// newline trimming
	tmpl = mustParse(t, "{{ 23 -}}\n<\r\n{{- 45 }}")
	result = mustExecute(t, tmpl, nil)
	if g, w := result, "23<45"; g != w {
		t.Errorf("got '%v', want '%v'", g, w)
	}
}

func TestActions(t *testing.T) {

	t.Run("comments", func(t *testing.T) {
		// comments
		tmpl := mustParse(t, " {{/* comment */}} ")
		result := mustExecute(t, tmpl, nil)
		if g, w := result, "  "; g != w {
			t.Errorf("got '%v', want '%v'", g, w)
		}

		// with trimming
		tmpl = mustParse(t, " {{- /* comment */ -}} ")
		result = mustExecute(t, tmpl, nil)
		if g, w := result, ""; g != w {
			t.Errorf("got '%v', want '%v'", g, w)
		}
	})

	t.Run("if", func(t *testing.T) {
		tmpl := mustParse(t, "{{- if .Message -}} present! {{- else -}} empty! {{- end -}}")
		result := mustExecute(t, tmpl, map[string]string{"Message": "foo"})
		if g, w := result, "present!"; g != w {
			t.Errorf("got '%v', want '%v'", g, w)
		}

		result = mustExecute(t, tmpl, map[string]string{"Message": ""})
		if g, w := result, "empty!"; g != w {
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
		tmpl := mustParse(t, ts)
		result := mustExecute(t, tmpl, map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		})

		want := `
Key:key1 Value:value1
Key:key2 Value:value2
Key:key3 Value:value3
`
		if g, w := strings.TrimSpace(result), strings.TrimSpace(want); g != w {
			t.Errorf("got '%v', want '%v'", g, w)
		}
	})

	t.Run("with", func(t *testing.T) {
		ts := strings.TrimSpace(`
{{with .Message}}
print {{.}}
{{end}}
`)
		tmpl := mustParse(t, ts)
		result := mustExecute(t, tmpl, map[string]string{
			"Message": "Hello",
		})

		if g, w := strings.TrimSpace(result), "print Hello"; g != w {
			t.Errorf("got '%v', want '%v'", g, w)
		}
	})
}

// MyStruct struct
type MyStruct struct {
}

// https://golang.org/pkg/text/template/#hdr-Arguments
// method must have one return value (of any type) or two return values, the second of which is an error.
func (s *MyStruct) Do() string {
	return "do!!!!"
}

func TestArguments(t *testing.T) {
	t.Run("nameless variable", func(t *testing.T) {
		tmpl := mustParse(t, "{{ $ := 100}}{{ $ }}")
		result := mustExecute(t, tmpl, nil)
		if g, w := result, "100"; g != w {
			t.Errorf("got %v, want %v", g, w)
		}
	})

	t.Run("invoke method", func(t *testing.T) {
		tmpl := mustParse(t, "{{ .Do }}")
		result := mustExecute(t, tmpl, &MyStruct{})
		if g, w := result, "do!!!!"; g != w {
			t.Errorf("got %v, want %v", g, w)
		}
	})

	t.Run("invoke function", func(t *testing.T) {
		fm := template.FuncMap{
			"unix": func() int { return int(time.Now().Unix()) },
		}
		tmpl, _ := template.New("test").Funcs(fm).Parse("{{ unix }}")
		result := mustExecute(t, tmpl, nil)
		log.Println(result)
	})
}

func TestPipelines(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		tmpl := mustParse(t, "{{ .foo | println }}")
		result := mustExecute(t, tmpl, map[string]string{"foo": "bar"})
		if g, w := result, "bar\n"; g != w {
			t.Errorf("got %v, want %v", g, w)
		}
	})
}

func mustParse(t *testing.T, text string) *template.Template {
	t.Helper()
	tmpl := template.New("test")
	tmpl, err := tmpl.Parse(text)
	if err != nil {
		t.Fatal(err)
	}
	return tmpl
}

func mustExecute(t *testing.T, tmpl *template.Template, data interface{}) string {
	t.Helper()
	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

package template

import (
	"bytes"
	"html/template"
)

// Unescape func
func Unescape() string {
	tmpl := template.New("foo")
	tmpl, err := tmpl.Parse("<html><body>{{ .Body }}</body></html>")
	if err != nil {
		panic(err)
	}
	var b bytes.Buffer
	data := map[string]interface{}{
		"Body": template.HTML("<p>Hello</p>"), // template.HTML で渡すことでエスケープ処理をスキップできる
	}
	if err := tmpl.Execute(&b, data); err != nil {
		panic(err)
	}
	return b.String()
}

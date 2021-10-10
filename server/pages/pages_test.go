package pages

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
)

func TestTemplates(t *testing.T) {
	var tmpl = template.Must(template.New("test").Parse(`dolbaeb {{if .loggedIn}}yes{{else}}no{{end}}`))

	var params = map[string]bool{
		"loggedIn": true,
	}
	buf := bytes.Buffer{}
	err := tmpl.Lookup("test").Execute(&buf, params)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(buf.String())
}

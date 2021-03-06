package main

import (
	"bytes"
	"log"
	"text/template"
)

const FileTemplate = `// Code generated by "{{.Flags}}"; DO NOT EDIT.
package {{.Package}}

{{range $k,$v := .Funcs}}
{{$v}} {{end}}
`

type FileTemplateContent struct {
	Flags   string
	Package string
	Funcs   []string
}

func (c FileTemplateContent) generateContent() ([]byte, error) {
	// file content container
	b := bytes.NewBuffer([]byte{})

	// new file template
	t, err := template.New("").Parse(FileTemplate)
	if err != nil {
		log.Fatal(err)
	}
	// generate file content write to file content container
	if err := t.Execute(b, c); err != nil {
		log.Fatal(err)
	}

	return b.Bytes(), nil
}

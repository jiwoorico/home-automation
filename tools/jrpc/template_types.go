package main

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jakewright/home-automation/tools/jrpc/imports"
)

const packageDirExternal = "def"

type typesDataMessage struct {
	Name   string
	Fields []*typesDataField
}

type typesDataField struct {
	GoName        string
	JSONName      string
	Type          string
	IsMessageType bool
	Repeated      bool
	Pointer       bool
	Required      bool
}

type typesData struct {
	PackageName string
	Imports     []*imports.Imp
	Messages    []*typesDataMessage
}

const typesTemplateText = `// Code generated by jrpc. DO NOT EDIT.

package {{ .PackageName }}

{{ if .Imports }}
	import (
		{{- range .Imports }}
			{{ .Alias }} "{{ .Path }}"
		{{- end}}
	)
{{- end }}

{{ range $message := .Messages }}
	// {{ $message.Name }} is defined in the .def file
	type {{ $message.Name }} struct {
		{{- range $field := .Fields }}
			{{ $field.GoName }} {{ $field.Type }} ` + "`" + `json:"{{ $field.JSONName }}"` + "`" + `
		{{- end }}
	}
{{- end }}

{{ range $message := .Messages }}
	// Validate returns an error if any of the fields have bad values
	func (m *{{ $message.Name }}) Validate() error {
		{{- range $field := $message.Fields -}}
			{{ if $field.IsMessageType -}}
				{{ if $field.Repeated -}}
					{{ if $field.Pointer -}}
						if m.{{ $field.GoName }} != nil {
							for _, r := range *m.{{ $field.GoName }} {
								if err := r.Validate(); err != nil {
									return err
								}
							}
						}
					{{ else -}}
						for _, r := range m.{{ $field.GoName }} {
							if err := r.Validate(); err != nil {
								return err
							}
						}
					{{ end -}}
				{{ else -}}
					if err := m.{{ $field.GoName }}.Validate(); err != nil {
						return err
					}
				{{ end }}
			{{ end -}}

			{{ if $field.Required -}}
				{{ if $field.Pointer -}}
					if m.{{ $field.GoName }} == nil {
				{{ else if $field.Repeated -}}
					if len(m.{{ $field.GoName }}) == 0 {
				{{ else if eq $field.Type "[]byte" -}}
					if len(m.{{ $field.GoName }}) == 0 {
				{{ else if eq $field.Type "string" -}}
					if m.{{ $field.GoName }} == "" {
				{{ else if eq $field.Type "int32" "int64" "uint32" "uint64" "float32" "float64" -}}
					if m.{{ $field.GoName }} == 0 {
				{{ else if eq $field.Type "time.Time" -}} 
					if m.{{ $field.GoName }}.IsZero() {
				{{ else }}
					if true {
				{{ end -}}
					return errors.BadRequest("field {{ $field.JSONName }} is required")
				}
			{{ end -}}
		{{ end -}}

		return nil
	}
{{ end }}

`

type typesGenerator struct {
	baseGenerator
}

func (g *typesGenerator) Template() (*template.Template, error) {
	return template.New("types_template").Parse(typesTemplateText)
}

func (g *typesGenerator) PackageDir() string {
	return packageDirExternal
}

func (g *typesGenerator) Data(im *imports.Manager) (interface{}, error) {
	im.Add("github.com/jakewright/home-automation/libraries/go/errors")

	if len(g.file.Messages) == 0 {
		return nil, nil
	}

	var messages []*typesDataMessage
	for _, m := range g.file.FlatMessages {
		alias, parts := m.Lineage()
		if alias != "" {
			// Ignore any messages that are from imported files
			continue
		}

		name := strings.Join(parts, "_")
		if !reValidGoStructUnderscore.MatchString(name) {
			return nil, fmt.Errorf("invalid message name %s", name)
		}

		fields := make([]*typesDataField, len(m.Fields))
		for i, f := range m.Fields {
			goName, jsonName, err := convertFieldName(f.Name)
			if err != nil {
				return nil, err
			}

			typ, err := resolveTypeName(f.Type, g.file, im)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve type of field %q in message %q: %w", f.Name, m.Name, err)
			}

			var required bool
			if v, ok := f.Options["required"].(bool); ok {
				required = v
			}

			fields[i] = &typesDataField{
				GoName:        goName,
				JSONName:      jsonName,
				Type:          typ.FullTypeName,
				IsMessageType: typ.IsMessageType,
				Repeated:      typ.Repeated,
				Pointer:       typ.Pointer,
				Required:      required,
			}
		}

		messages = append(messages, &typesDataMessage{
			Name:   name,
			Fields: fields,
		})
	}

	return &typesData{
		PackageName: externalPackageName(g.options),
		Imports:     im.Get(),
		Messages:    messages,
	}, nil
}

func (g *typesGenerator) Filename() string {
	return "types.go"
}

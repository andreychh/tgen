// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"embed"
	"slices"
	"text/template"
)

// templates holds the embedded [text/template] files used for generating Go
// source code.
//
//go:embed templates/*.tmpl
var templates embed.FS

type Template struct{}

func NewTemplate() Template {
	return Template{}
}

func (t Template) Value() (*template.Template, error) {
	return template.New("").
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"objects":                slices.Collect[Object],
			"methods":                slices.Collect[Method],
			"fields":                 slices.Collect[Field],
			"discriminated_unions":   slices.Collect[DiscriminatedUnion],
			"discriminated_variants": slices.Collect[DiscriminatedVariant],
			"discriminated_objects":  slices.Collect[DiscriminatedObject],
		}).
		ParseFS(templates, "templates/*.tmpl")
}

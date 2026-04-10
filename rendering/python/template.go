// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"embed"
	"slices"
	"text/template"
)

// templates holds the embedded [text/template] files used for generating Python
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
			"objects":               slices.Collect[Object],
			"discriminated_objects": slices.Collect[DiscriminatedObject],
			"discriminated_unions":  slices.Collect[DiscriminatedUnion],
			"fields":                slices.Collect[Field],
		}).
		ParseFS(templates, "templates/*.tmpl")
}

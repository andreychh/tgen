// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package output

import (
	"embed"
	"text/template"
)

// Mold represents a set of embedded templates configured with a function map,
// ready to produce a [template.Template].
type Mold struct {
	fs    embed.FS
	funcs template.FuncMap
}

// NewMold creates a Mold from the given embedded filesystem and function map.
func NewMold(fs embed.FS, funcs template.FuncMap) Mold {
	return Mold{fs: fs, funcs: funcs}
}

// Template parses all templates in the embedded filesystem and returns a
// configured [template.Template].
func (m Mold) Template() (*template.Template, error) {
	return template.New("").
		Option("missingkey=error").
		Funcs(m.funcs).
		ParseFS(m.fs, "templates/*.tmpl")
}

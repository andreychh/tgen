// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package output

import (
	"embed"
	"errors"
	"fmt"
	"strings"
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
// configured [template.Template]. Beside the configured function map, the
// returned template exposes three functions: render, which executes the
// template named by its first argument against its second and yields the
// rendered text; override, which names the template rendering one definition —
// the block written by hand for the reference it takes when the templates hold
// one, and the shape it takes otherwise; and assert, which stops the rendering
// when what a block was written for stops holding.
func (m Mold) Template() (*template.Template, error) {
	tmpl := template.New("").Option("missingkey=error").Funcs(m.funcs)
	return tmpl.Funcs(template.FuncMap{
		// assert states what a block written by hand takes for granted, and fails the
		// rendering with the given message when the specification stops granting it. A
		// block cannot be generated, so a requirement it does not already answer cannot
		// be answered by leaving something out — only by a person reading the block
		// again, which is what a stopped generation is for.
		"assert": func(holds bool, message string) (string, error) {
			if !holds {
				return "", errors.New(message)
			}
			return "", nil
		},
		"render": func(name string, data any) (string, error) {
			var b strings.Builder
			err := tmpl.ExecuteTemplate(&b, name, data)
			if err != nil {
				return "", fmt.Errorf("rendering %q: %w", name, err)
			}
			return b.String(), nil
		},
		// "manual_" prefixes the name of a template written by hand for one definition,
		// holding such names in a namespace of their own. A definition is addressed by
		// a reference the documentation gives it, a shape by a name a target chooses;
		// the prefix keeps a target from ever choosing a name a reference already
		// answers to, which would take a declaration over without saying so anywhere.
		"override": func(ref, shape string) string {
			name := "manual_" + ref
			if tmpl.Lookup(name) == nil {
				return shape
			}
			return name
		},
	}).ParseFS(m.fs, "templates/*.tmpl")
}

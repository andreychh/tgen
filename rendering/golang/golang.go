// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package golang provides the necessary templates and execution context to
// generate Go source code from the parsed Telegram Bot API specification.
package golang

import (
	"embed"
	"slices"
	"text/template"

	"github.com/andreychh/tgen/parsing"
)

// templates holds the embedded [text/template] files used for generating Go
// source code.
//
//go:embed templates/*.tmpl
var templates embed.FS

// PrepareTemplate initializes and parses the embedded Go code templates. It
// panics if the embedded templates contain syntax errors.
func PrepareTemplate() *template.Template {
	return template.Must(
		template.New("").Option("missingkey=error").Funcs(
			template.FuncMap{
				"objects":  slices.Collect[parsing.Object],
				"fields":   slices.Collect[parsing.Field],
				"unions":   slices.Collect[parsing.Union],
				"variants": slices.Collect[parsing.Variant],
			},
		).ParseFS(templates, "templates/*.tmpl"),
	)
}

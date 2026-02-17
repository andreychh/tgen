// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"embed"
	"slices"
	"text/template"

	"github.com/andreychh/tgen/parsing"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

func PrepareTemplate() *template.Template {
	tmpl := template.New("base").
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"unions":   slices.Collect[parsing.Union],
			"variants": slices.Collect[parsing.Variant],
		})
	return template.Must(tmpl.ParseFS(templateFS, "templates/*.tmpl"))
}

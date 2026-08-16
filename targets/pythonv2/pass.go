// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"embed"
	"fmt"
	"text/template"

	"github.com/andreychh/tgen/output"
)

//go:embed templates/*.tmpl
var templates embed.FS

// Pass is the Python generation stage: it renders the records of the
// pipeline's exit into the files of a Python package.
type Pass struct {
	gen Generation
}

// NewPass creates a Pass rendering the given generation.
func NewPass(gen Generation) Pass {
	return Pass{gen: gen}
}

// Artifacts returns the files the target writes, each bound to the template
// rendering it. It fails when a template is malformed.
func (p Pass) Artifacts() (output.Artifacts, error) {
	tmpl, err := output.NewMold(templates, template.FuncMap{}).Template()
	if err != nil {
		return nil, fmt.Errorf("preparing template: %w", err)
	}
	return output.Artifacts{
		"api.py": output.NewTemplateView(tmpl, "api", p.gen),
	}, nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"io"
	"text/template"

	"github.com/andreychh/tgen/parsing"
)

type Unions struct {
	spec     parsing.Specification
	template *template.Template
}

func NewUnions(spec parsing.Specification, tmpl *template.Template) Unions {
	return Unions{
		spec:     spec,
		template: tmpl,
	}
}

func (u Unions) Render(w io.Writer) error {
	return u.template.ExecuteTemplate(w, "unions", u.spec)
}

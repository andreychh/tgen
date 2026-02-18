// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package rendering

import (
	"io"
	"text/template"
)

// TemplateView implements the [View] interface by executing a specific
// [text/template] with the provided data.
type TemplateView struct {
	template *template.Template
	name     string
	data     any
}

// NewTemplateView returns a [TemplateView] configured to execute the template
// named name using the provided data.
func NewTemplateView(t *template.Template, name string, data any) TemplateView {
	return TemplateView{
		template: t,
		name:     name,
		data:     data,
	}
}

// Render executes the underlying template and writes the generated output to w.
func (v TemplateView) Render(w io.Writer) error {
	return v.template.ExecuteTemplate(w, v.name, v.data)
}

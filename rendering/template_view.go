// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package rendering

import (
	"io"
	"text/template"
)

// TemplateView represents a [View] that renders a named [text/template] with
// structured data.
type TemplateView struct {
	template *template.Template
	name     string
	data     any
}

// NewTemplateView constructs a [TemplateView] for template name in t with the
// given data.
func NewTemplateView(t *template.Template, name string, data any) TemplateView {
	return TemplateView{
		template: t,
		name:     name,
		data:     data,
	}
}

// Render implements [View].
func (v TemplateView) Render(w io.Writer) error {
	return v.template.ExecuteTemplate(w, v.name, v.data)
}

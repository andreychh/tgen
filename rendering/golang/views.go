// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"text/template"

	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/rendering"
)

// NewObjectsView returns a [rendering.TemplateView] configured to render the
// "objects" template using the provided [parsing.Specification].
func NewObjectsView(t *template.Template, s parsing.Specification) rendering.TemplateView {
	return rendering.NewTemplateView(t, "objects", s)
}

// NewUnionsView returns a [rendering.TemplateView] configured to render the
// "unions" template using the provided [parsing.Specification].
func NewUnionsView(t *template.Template, s parsing.Specification) rendering.TemplateView {
	return rendering.NewTemplateView(t, "unions", s)
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"embed"
	"fmt"
	"slices"
	"text/template"

	"github.com/andreychh/tgen/meta"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/output"
	"github.com/andreychh/tgen/targets"
)

//go:embed templates/*.tmpl
var templates embed.FS

// Pass assembles the output artifacts for the Python code generation target.
type Pass struct {
	context GenerationContext
}

// NewPass creates a Pass for the given specification and snapshot.
func NewPass(spec explicit.Specification, snapshot meta.Snapshot) Pass {
	return Pass{context: NewGenerationContext(
		NewSpecification(spec),
		targets.NewSnapshot(snapshot),
	)}
}

// Artifacts produces the output artifacts for the Python code generation target.
func (p Pass) Artifacts() (output.Artifacts, error) {
	mold := output.NewMold(templates, template.FuncMap{
		"objects":               slices.Collect[Object],
		"discriminated_objects": slices.Collect[DiscriminatedObject],
		"discriminated_unions":  slices.Collect[DiscriminatedUnion],
		"methods":               slices.Collect[Method],
		"fields":                slices.Collect[Field],
	})
	tmpl, err := mold.Template()
	if err != nil {
		return nil, fmt.Errorf("preparing template: %w", err)
	}
	return output.Artifacts{
		"__init__.py":         output.NewTemplateView(tmpl, "init", p.context),
		"types.py":            output.NewTemplateView(tmpl, "types", p.context),
		"method.py":           output.NewTemplateView(tmpl, "method_enum", p.context),
		"methods.py":          output.NewTemplateView(tmpl, "methods", p.context),
		"client.py":           output.NewTemplateView(tmpl, "client", p.context),
		"asyncio/__init__.py": output.NewTemplateView(tmpl, "asyncio_init", p.context),
		"asyncio/methods.py":  output.NewTemplateView(tmpl, "async_methods", p.context),
		"asyncio/client.py":   output.NewTemplateView(tmpl, "async_client", p.context),
		"asyncio/fake.py":     output.NewTemplateView(tmpl, "async_fake", p.context),
	}, nil
}

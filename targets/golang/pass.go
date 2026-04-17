// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

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

// Pass assembles the output artifacts for the Go code generation target.
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

// Artifacts produces the output artifacts for the Go code generation target.
func (p Pass) Artifacts() (output.Artifacts, error) {
	mold := output.NewMold(templates, template.FuncMap{
		"objects":                slices.Collect[Object],
		"methods":                slices.Collect[Method],
		"fields":                 slices.Collect[Field],
		"discriminated_unions":   slices.Collect[DiscriminatedUnion],
		"discriminated_variants": slices.Collect[DiscriminatedVariant],
		"discriminated_objects":  slices.Collect[DiscriminatedObject],
	})
	tmpl, err := mold.Template()
	if err != nil {
		return nil, fmt.Errorf("preparing template: %w", err)
	}
	return output.Artifacts{
		"objects.go": output.NewTemplateView(tmpl, "objects", p.context),
		"unions.go":  output.NewTemplateView(tmpl, "unions", p.context),
		"methods.go": output.NewTemplateView(tmpl, "methods", p.context),
		"client.go":  output.NewTemplateView(tmpl, "client", p.context),
	}, nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/enrichment"
	"github.com/andreychh/tgen/meta"
	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/rendering"
)

// Artifacts assembles the rendering artifacts for the Go code generation target.
type Artifacts struct {
	spec     parsing.Specification
	snapshot meta.Snapshot
}

// NewArtifacts creates an Artifacts for the given specification and snapshot.
func NewArtifacts(spec parsing.Specification, snapshot meta.Snapshot) Artifacts {
	return Artifacts{spec: spec, snapshot: snapshot}
}

func (a Artifacts) Value() (rendering.Artifacts, error) {
	tmpl, err := NewTemplate().Value()
	if err != nil {
		return nil, fmt.Errorf("preparing template: %w", err)
	}
	ctx := NewGenerationContext(
		NewSpecification(enrichment.NewSpecification(a.spec)),
		rendering.NewSnapshot(a.snapshot),
	)
	return rendering.Artifacts{
		"objects.go": rendering.NewTemplateView(tmpl, "objects", ctx),
		"unions.go":  rendering.NewTemplateView(tmpl, "unions", ctx),
		"methods.go": rendering.NewTemplateView(tmpl, "methods", ctx),
	}, nil
}

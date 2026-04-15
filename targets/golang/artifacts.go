// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/meta"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/output"
)

// Artifacts assembles the rendering artifacts for the Go code generation target.
type Artifacts struct {
	spec     explicit.Specification
	snapshot meta.Snapshot
}

// NewArtifacts creates an Artifacts for the given specification and snapshot.
func NewArtifacts(spec explicit.Specification, snapshot meta.Snapshot) Artifacts {
	return Artifacts{spec: spec, snapshot: snapshot}
}

func (a Artifacts) Value() (output.Artifacts, error) {
	tmpl, err := NewTemplate().Value()
	if err != nil {
		return nil, fmt.Errorf("preparing template: %w", err)
	}
	ctx := NewGenerationContext(NewSpecification(a.spec), output.NewSnapshot(a.snapshot))
	return output.Artifacts{
		"objects.go": output.NewTemplateView(tmpl, "objects", ctx),
		"unions.go":  output.NewTemplateView(tmpl, "unions", ctx),
		"methods.go": output.NewTemplateView(tmpl, "methods", ctx),
		"client.go":  output.NewTemplateView(tmpl, "client", ctx),
		"fake.go":    output.NewTemplateView(tmpl, "fake", ctx),
	}, nil
}

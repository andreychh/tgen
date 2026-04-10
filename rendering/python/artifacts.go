// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"fmt"

	"github.com/andreychh/tgen/meta"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/rendering"
)

// Artifacts assembles the rendering artifacts for the Python code generation target.
type Artifacts struct {
	spec     explicit.Specification
	snapshot meta.Snapshot
}

// NewArtifacts creates an Artifacts for the given specification and snapshot.
func NewArtifacts(spec explicit.Specification, snapshot meta.Snapshot) Artifacts {
	return Artifacts{spec: spec, snapshot: snapshot}
}

func (a Artifacts) Value() (rendering.Artifacts, error) {
	tmpl, err := NewTemplate().Value()
	if err != nil {
		return nil, fmt.Errorf("preparing template: %w", err)
	}
	ctx := NewGenerationContext(NewSpecification(a.spec), rendering.NewSnapshot(a.snapshot))
	return rendering.Artifacts{
		"objects.py": rendering.NewTemplateView(tmpl, "objects", ctx),
	}, nil
}

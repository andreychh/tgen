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
		"__init__.py":         rendering.NewTemplateView(tmpl, "init", ctx),
		"types.py":            rendering.NewTemplateView(tmpl, "types", ctx),
		"method.py":           rendering.NewTemplateView(tmpl, "method_enum", ctx),
		"methods.py":          rendering.NewTemplateView(tmpl, "methods", ctx),
		"client.py":           rendering.NewTemplateView(tmpl, "client", ctx),
		"fake.py":             rendering.NewTemplateView(tmpl, "fake", ctx),
		"asyncio/__init__.py": rendering.NewTemplateView(tmpl, "asyncio_init", ctx),
		"asyncio/methods.py":  rendering.NewTemplateView(tmpl, "async_methods", ctx),
		"asyncio/client.py":   rendering.NewTemplateView(tmpl, "async_client", ctx),
		"asyncio/fake.py":     rendering.NewTemplateView(tmpl, "async_fake", ctx),
	}, nil
}

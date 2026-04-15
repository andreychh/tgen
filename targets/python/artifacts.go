// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"fmt"

	"github.com/andreychh/tgen/meta"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/output"
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

func (a Artifacts) Value() (output.Artifacts, error) {
	tmpl, err := NewTemplate().Value()
	if err != nil {
		return nil, fmt.Errorf("preparing template: %w", err)
	}
	ctx := NewGenerationContext(NewSpecification(a.spec), output.NewSnapshot(a.snapshot))
	return output.Artifacts{
		"__init__.py":         output.NewTemplateView(tmpl, "init", ctx),
		"types.py":            output.NewTemplateView(tmpl, "types", ctx),
		"method.py":           output.NewTemplateView(tmpl, "method_enum", ctx),
		"methods.py":          output.NewTemplateView(tmpl, "methods", ctx),
		"client.py":           output.NewTemplateView(tmpl, "client", ctx),
		"fake.py":             output.NewTemplateView(tmpl, "fake", ctx),
		"asyncio/__init__.py": output.NewTemplateView(tmpl, "asyncio_init", ctx),
		"asyncio/methods.py":  output.NewTemplateView(tmpl, "async_methods", ctx),
		"asyncio/client.py":   output.NewTemplateView(tmpl, "async_client", ctx),
		"asyncio/fake.py":     output.NewTemplateView(tmpl, "async_fake", ctx),
	}, nil
}

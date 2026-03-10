// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/rendering"
)

type Artifacts struct {
	spec parsing.Specification
}

func NewArtifacts(spec parsing.Specification) Artifacts {
	return Artifacts{spec: spec}
}

func (a Artifacts) Value() (rendering.Artifacts, error) {
	tmpl, err := NewTemplate().Value()
	if err != nil {
		return nil, fmt.Errorf("preparing template: %w", err)
	}
	spec := NewSpecification(a.spec)
	return rendering.Artifacts{
		"objects.go": rendering.NewTemplateView(tmpl, "objects", spec),
		"unions.go":  rendering.NewTemplateView(tmpl, "unions", spec),
	}, nil
}

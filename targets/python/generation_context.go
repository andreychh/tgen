// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"github.com/andreychh/tgen/targets"
)

// GenerationContext holds all data required to render a single code generation run.
type GenerationContext struct {
	spec     Specification
	snapshot targets.Snapshot
}

func NewGenerationContext(spec Specification, snapshot targets.Snapshot) GenerationContext {
	return GenerationContext{
		spec:     spec,
		snapshot: snapshot,
	}
}

func (c GenerationContext) Spec() Specification {
	return c.spec
}

func (c GenerationContext) Snapshot() targets.Snapshot {
	return c.snapshot
}

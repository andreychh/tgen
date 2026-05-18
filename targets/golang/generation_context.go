// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"github.com/andreychh/tgen/targets"
)

// GenerationContext holds all data required to render a single code generation run.
type GenerationContext struct {
	spec        Specification
	snapshot    targets.Snapshot
	packageName string
}

func NewGenerationContext(spec Specification, snapshot targets.Snapshot) GenerationContext {
	return GenerationContext{
		spec:        spec,
		snapshot:    snapshot,
		packageName: "api",
	}
}

func (c GenerationContext) Spec() Specification {
	return c.spec
}

func (c GenerationContext) Snapshot() targets.Snapshot {
	return c.snapshot
}

func (c GenerationContext) PackageName() string {
	return c.packageName
}

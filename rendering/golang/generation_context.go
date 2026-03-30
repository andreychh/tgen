// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"github.com/andreychh/tgen/rendering"
)

// GenerationContext holds all data required to render a single code generation run.
type GenerationContext struct {
	spec     Specification
	snapshot rendering.Snapshot
}

func NewGenerationContext(spec Specification, snapshot rendering.Snapshot) GenerationContext {
	return GenerationContext{
		spec:     spec,
		snapshot: snapshot,
	}
}

func (c GenerationContext) Spec() Specification { //nolint:gocritic // value receiver is intentional; GenerationContext is immutable
	return c.spec
}

func (c GenerationContext) Snapshot() rendering.Snapshot { //nolint:gocritic // value receiver is intentional; GenerationContext is immutable
	return c.snapshot
}

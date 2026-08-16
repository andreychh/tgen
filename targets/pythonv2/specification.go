// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package pythonv2 renders the records of the pipeline's exit as Python source:
// it spells every type the way a pydantic model spells it and picks, for each
// method, the template that assembles its request and the one that names what
// it hands back.
package pythonv2

import (
	"fmt"

	ir "github.com/andreychh/tgen/model/ir/v2"
	"github.com/andreychh/tgen/pkg/slices"
)

// Specification represents the Python view of the specification: the
// declarations each generated file is rendered from, and the release those
// declarations were read from. What a run decides rather than the documentation
// belongs to [Generation].
type Specification struct {
	inner ir.Specification
}

// NewSpecification creates a Specification over the records of the pipeline's
// exit.
func NewSpecification(inner ir.Specification) Specification {
	return Specification{inner: inner}
}

// Release returns the Bot API release the specification was read from.
func (s Specification) Release() Release {
	return NewRelease(s.inner.Release())
}

// Definitions returns the declarations the generated package holds, ordered by
// the position the source of each record gave it. It fails when a record cannot
// be read as the declaration it is rendered as.
func (s Specification) Definitions() ([]Declaration, error) {
	records, err := s.inner.Definitions()
	if err != nil {
		return nil, fmt.Errorf("reading definitions: %w", err)
	}
	return slices.NewMapped(records, NewDeclaration), nil
}

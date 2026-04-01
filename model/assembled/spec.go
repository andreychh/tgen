// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package assembled assembles the complete specification consumed by the
// rendering layer, combining the overlay-corrected explicit spec with the
// tgen-invented implicit types. It is the single entry point that rendering and
// the CLI should depend on.
package assembled

import (
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/model/explicit/overlays"
	"github.com/andreychh/tgen/model/implicit"
	"github.com/andreychh/tgen/model/implicit/literals"
)

// Specification represents the fully assembled Telegram Bot API specification,
// combining the overlay-modified explicit spec with tgen-defined implicit
// types.
type Specification struct {
	explicit explicit.Specification
	implicit implicit.Specification
}

// NewSpec constructs a Specification from a parsed specification.
func NewSpec(s explicit.Specification) Specification {
	return Specification{
		explicit: overlays.NewSpecification(s),
		implicit: literals.NewSpecification(),
	}
}

// Explicit returns the overlay-modified view of the Telegram Bot API
// specification.
func (s Specification) Explicit() explicit.Specification {
	return s.explicit
}

// Implicit returns the set of tgen-defined types not present in the
// specification.
func (s Specification) Implicit() implicit.Specification {
	return s.implicit
}

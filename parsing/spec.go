// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/parsing/explicit/overlays"
	"github.com/andreychh/tgen/parsing/implicit"
)

// Spec represents the fully enriched Telegram Bot API specification, combining
// the overlay-modified explicit spec with tgen-defined implicit types.
type Spec struct {
	explicit explicit.Specification
	implicit implicit.Spec
}

// NewSpec constructs a Spec from a parsed specification.
func NewSpec(s explicit.Specification) Spec {
	return Spec{
		explicit: overlays.NewSpec(s),
		implicit: implicit.NewSpec(),
	}
}

// Explicit returns the overlay-modified view of the Telegram Bot API
// specification.
func (s Spec) Explicit() explicit.Specification { return s.explicit }

// Implicit returns the set of tgen-defined types not present in the
// specification.
func (s Spec) Implicit() implicit.Spec { return s.implicit }

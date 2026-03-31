// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// Unions represents an enriched union collection with overlay applied to all
// variant fields.
type Unions struct {
	inner   explicit.Unions
	overlay Overlay
}

// NewUnions constructs a Unions from parsed unions with the given overlay
// applied to variant fields.
func NewUnions(u explicit.Unions, o Overlay) Unions {
	return Unions{inner: u, overlay: o}
}

func (u Unions) Discriminated() iter.Seq[explicit.DiscriminatedUnion] {
	return iters.NewMappedSeq(
		u.inner.Discriminated(),
		func(d explicit.DiscriminatedUnion) explicit.DiscriminatedUnion {
			return NewDiscriminatedUnion(d, u.overlay)
		},
	)
}

func (u Unions) Structured() iter.Seq[explicit.StructuredUnion] {
	return iters.NewMappedSeq(
		u.inner.Structured(),
		func(s explicit.StructuredUnion) explicit.StructuredUnion {
			return NewStructuredUnion(s, u.overlay)
		},
	)
}

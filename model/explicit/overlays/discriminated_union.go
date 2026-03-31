// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// DiscriminatedUnion represents an enriched discriminated union with overlay
// applied to variant fields.
type DiscriminatedUnion struct {
	inner   explicit.DiscriminatedUnion
	overlay Overlay
}

// NewDiscriminatedUnion constructs a DiscriminatedUnion from a parsed union
// with the given overlay applied to variant fields.
func NewDiscriminatedUnion(u explicit.DiscriminatedUnion, o Overlay) DiscriminatedUnion {
	return DiscriminatedUnion{inner: u, overlay: o}
}

func (u DiscriminatedUnion) Reference() model.Reference {
	return u.inner.Reference()
}

func (u DiscriminatedUnion) Name() model.Name {
	return u.inner.Name()
}

func (u DiscriminatedUnion) Description() model.Description {
	return u.inner.Description()
}

func (u DiscriminatedUnion) DiscriminatorKey() model.Key {
	return u.inner.DiscriminatorKey()
}

func (u DiscriminatedUnion) Variants() iter.Seq[explicit.DiscriminatedVariant] {
	return iters.NewMappedSeq(
		u.inner.Variants(),
		func(v explicit.DiscriminatedVariant) explicit.DiscriminatedVariant {
			return NewDiscriminatedVariant(v, u.overlay)
		},
	)
}

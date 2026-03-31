// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// StructuredUnion represents an enriched structured union with overlay applied
// to variant fields.
type StructuredUnion struct {
	inner   explicit.StructuredUnion
	overlay Overlay
}

// NewStructuredUnion constructs a StructuredUnion from a parsed union with the
// given overlay applied to variant fields.
func NewStructuredUnion(u explicit.StructuredUnion, o Overlay) StructuredUnion {
	return StructuredUnion{inner: u, overlay: o}
}

func (u StructuredUnion) Reference() model.Reference {
	return u.inner.Reference()
}

func (u StructuredUnion) Name() model.Name {
	return u.inner.Name()
}

func (u StructuredUnion) Description() model.Description {
	return u.inner.Description()
}

func (u StructuredUnion) Variants() iter.Seq[explicit.StructuredVariant] {
	return iters.NewMappedSeq(
		u.inner.Variants(),
		func(v explicit.StructuredVariant) explicit.StructuredVariant {
			return NewStructuredVariant(v, u.overlay)
		},
	)
}

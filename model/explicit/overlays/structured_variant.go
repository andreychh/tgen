// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// StructuredVariant represents an enriched structured variant with overlay
// applied to its fields.
type StructuredVariant struct {
	inner   explicit.StructuredVariant
	overlay Overlay
}

// NewStructuredVariant constructs a StructuredVariant from a parsed variant
// with the given overlay applied to its fields.
func NewStructuredVariant(v explicit.StructuredVariant, o Overlay) StructuredVariant {
	return StructuredVariant{inner: v, overlay: o}
}

func (v StructuredVariant) Reference() model.Reference {
	return v.inner.Reference()
}

func (v StructuredVariant) Name() model.Name {
	return v.inner.Name()
}

func (v StructuredVariant) Description() model.Description {
	return v.inner.Description()
}

func (v StructuredVariant) Fields() iter.Seq[explicit.Field] {
	return iters.NewMappedSeq(v.inner.Fields(), v.overlay.Apply)
}

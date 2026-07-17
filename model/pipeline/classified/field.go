// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package classified

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/typed"
)

// Field is a classified field, an alias of [typed.Field]: the classified
// stage changes which fields belong to Fields, not the shape of any single
// field.
type Field = typed.Field

// FieldFilter is a [pipeline.Filter] that excludes a field already extracted
// into discriminators.
type FieldFilter struct {
	discriminators Discriminators
}

// NewFieldFilter constructs a FieldFilter over the discriminators already
// extracted from Fields.
func NewFieldFilter(discriminators Discriminators) FieldFilter {
	return FieldFilter{discriminators: discriminators}
}

// Apply implements [pipeline.Filter]. It reports whether field belongs in
// the filtered result: false when its owner already carries a discriminator
// under field's own key.
func (f FieldFilter) Apply(key model.FieldKey, field Field) bool {
	record, exists := f.discriminators.Lookup(key.Owner)
	return !exists || record.Key != key.Key
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package classified

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/classified/discriminator"
	"github.com/andreychh/tgen/model/pipeline/typed"
)

// Discriminator is the fixed value a field's description decodes for the
// reference that owns it.
type Discriminator struct {
	// Key is the field that carries the fixed value.
	Key   model.Key
	Value model.DiscriminatorValue
}

// DiscriminatorTable is the extraction operator: it decodes every field at
// position zero for the fixed discriminator value its description may carry,
// keyed by the reference of the field's owner.
type DiscriminatorTable struct {
	fields typed.Fields
}

// NewDiscriminatorTable constructs a DiscriminatorTable over fields.
func NewDiscriminatorTable(fields typed.Fields) DiscriminatorTable {
	return DiscriminatorTable{fields: fields}
}

// Apply returns the discriminators table, one record per field at position
// zero whose description decodes a fixed discriminator value, keyed by the
// reference of the field's owner.
func (t DiscriminatorTable) Apply() Discriminators {
	out := pipeline.NewMapTable[model.Reference, Discriminator]()
	for key, field := range t.fields.All() {
		if field.Position != 0 {
			continue
		}
		value, ok := discriminator.NewLabel(field.Description).Value()
		if !ok {
			continue
		}
		out.Insert(key.Owner, Discriminator{Key: field.Key, Value: value})
	}
	return out
}

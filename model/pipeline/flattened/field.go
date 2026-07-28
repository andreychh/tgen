// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package flattened

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/typed"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeform"
)

// Field is a field of an object or a parameter of a method whose type is
// reduced to a flat one. Its key, position, optionality, and description carry
// over from the corrected stage.
type Field struct {
	Key         model.Key
	Position    model.Position
	Type        typeform.Type
	Optionality model.Optionality
	Description prose.Phrase
}

// FieldMapping maps a corrected field into a flattened field by reducing its
// type expression to a flat type.
type FieldMapping struct{}

// NewFieldMapping constructs a FieldMapping.
func NewFieldMapping() FieldMapping {
	return FieldMapping{}
}

// Apply implements [pipeline.Mapping]. It fails when the field's type holds a
// union.
func (m FieldMapping) Apply(field typed.Field) (Field, error) {
	typ, err := NewNarrowing(field.Type).Value()
	if err != nil {
		return Field{}, fmt.Errorf("narrowing field type: %w", err)
	}
	return Field{
		Key:         field.Key,
		Position:    field.Position,
		Type:        typ,
		Optionality: field.Optionality,
		Description: field.Description,
	}, nil
}

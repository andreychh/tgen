// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package typed

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/typed/types"
	"github.com/andreychh/tgen/model/pipeline/unified"
	"github.com/andreychh/tgen/model/prose"
	typetree "github.com/andreychh/tgen/model/types/v2"
)

// Field is a field of an object or a parameter of a method whose type is
// resolved from prose into a type expression. Its key, optionality, and
// description carry over from the unified stage.
type Field struct {
	Key         model.Key
	Type        typetree.Expression
	Optionality model.Optionality
	Description prose.Phrase
}

// FieldMapping maps a unified field into a typed field by decoding the prose of
// its type column into a type expression.
type FieldMapping struct{}

// NewFieldMapping constructs a FieldMapping.
func NewFieldMapping() FieldMapping {
	return FieldMapping{}
}

// Apply implements [pipeline.Mapping]. It fails when the type prose does not
// decode into a valid type expression.
func (m FieldMapping) Apply(field unified.Field) (Field, error) {
	expr, err := types.NewExpression(field.Type).Value()
	if err != nil {
		return Field{}, fmt.Errorf("decoding type expression: %w", err)
	}
	return Field{
		Key:         field.Key,
		Type:        expr,
		Optionality: field.Optionality,
		Description: field.Description,
	}, nil
}

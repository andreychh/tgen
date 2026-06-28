// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package types decodes the prose of a field's type column into a type
// expression.
package types

import (
	"fmt"

	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/types/v2"
)

// Expression is the prose of a field's type column, ready to be decoded into
// a type expression.
type Expression struct {
	phrase prose.Phrase
}

// NewExpression constructs an Expression over the prose of a type column.
func NewExpression(phrase prose.Phrase) Expression {
	return Expression{phrase: phrase}
}

// Value returns the type expression decoded from the prose. It fails when the
// prose lexes to an unknown token or does not form a valid type expression.
//
//nolint:ireturn // types.Expression is the intentional decoded contract.
func (e Expression) Value() (types.Expression, error) {
	tokens, err := NewLexer(e.phrase).Tokens()
	if err != nil {
		return nil, fmt.Errorf("lexing type expression: %w", err)
	}
	expr, err := NewParser(NewCursor(tokens)).Expression()
	if err != nil {
		return nil, fmt.Errorf("parsing type expression: %w", err)
	}
	return expr, nil
}

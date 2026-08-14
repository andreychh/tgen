// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package types decodes the prose of a field's type column into a type
// expression.
//
// [Expression] is the entry point: wrap the prose of a type cell and call its
// Value method to obtain the type. Decoding runs in two stages — [Lexer] reads
// the prose into [Token] values, and [Parser] reads those against the grammar
// of a type cell.
package types

import (
	"fmt"

	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeexpr"
)

// Expression is the prose of a field's type column, ready to be decoded into a
// type expression.
type Expression struct {
	phrase prose.Phrase
}

// NewExpression constructs an Expression over the prose of a type column.
func NewExpression(phrase prose.Phrase) Expression {
	return Expression{phrase: phrase}
}

// Value returns the type expression decoded from the prose. It fails when the
// prose holds a word or a link that cannot be lexed, or when the tokens do not
// form a single valid type expression.
func (e Expression) Value() (typeexpr.Expression, error) {
	tokens, err := NewLexer(e.phrase).Tokens()
	if err != nil {
		return nil, fmt.Errorf("lexing type expression: %w", err)
	}
	expr, err := NewParser(tokens).Expression()
	if err != nil {
		return nil, fmt.Errorf("parsing type expression: %w", err)
	}
	return expr, nil
}

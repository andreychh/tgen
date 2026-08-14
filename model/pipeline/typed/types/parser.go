// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import (
	"errors"

	"github.com/andreychh/tgen/model/typeexpr"
)

// Parser decodes a stream of [Token] values into the type expression they
// spell. The grammar it accepts is:
//
//	expression := term ("or" term)*
//	term       := "Array of" element | ref | primitive
//	element    := term (("," | "and") term)*
//
// The documentation writes two separators and they bind differently. An "or"
// joins whole terms, so "Array of A or B" is a union of an array and B; a comma
// list closed by "and" joins the members of an array's element union, so "Array
// of A, B and C" is an array of a union.
type Parser struct {
	tokens []Token
}

// NewParser constructs a Parser over the tokens of one type expression.
func NewParser(tokens []Token) Parser {
	return Parser{tokens: tokens}
}

// Expression returns the type expression decoded from the whole token stream.
// It fails when the tokens are malformed or do not all form a single
// expression.
func (p Parser) Expression() (typeexpr.Expression, error) {
	cursor := NewCursor(p.tokens)
	expr, err := p.expression(cursor)
	if err != nil {
		return nil, err
	}
	if !cursor.Done() {
		return nil, errors.New("trailing tokens after type expression")
	}
	return expr, nil
}

// expression parses the alternatives of a whole type, joined by "or".
func (p Parser) expression(cursor *Cursor[Token]) (typeexpr.Expression, error) {
	return p.union(cursor, NewOr())
}

// element parses the members of an array's element enumeration, joined by
// commas and a closing "and".
func (p Parser) element(cursor *Cursor[Token]) (typeexpr.Expression, error) {
	return p.union(cursor, NewSeries())
}

// union returns the [typeexpr.Union] of every term that sep joins, or the lone
// term when sep never follows. It fails when a term is malformed.
func (p Parser) union(cursor *Cursor[Token], sep Token) (typeexpr.Expression, error) {
	first, err := p.term(cursor)
	if err != nil {
		return nil, err
	}
	variants := []typeexpr.Expression{first}
	for p.at(cursor, sep) {
		cursor.Skip()
		next, err := p.term(cursor)
		if err != nil {
			return nil, err
		}
		variants = append(variants, next)
	}
	if len(variants) == 1 {
		return variants[0], nil
	}
	return typeexpr.NewUnion(variants...), nil
}

// term parses an "Array of" wrapper around an element enumeration, or a single
// reference or primitive. It fails on the end of input or a token that cannot
// begin a type.
func (p Parser) term(cursor *Cursor[Token]) (typeexpr.Expression, error) {
	node, ok := cursor.Take()
	if !ok {
		return nil, errors.New("expected a type, found end of input")
	}
	switch node := node.(type) {
	case ArrayOf:
		element, err := p.element(cursor)
		if err != nil {
			return nil, err
		}
		return typeexpr.NewArray(element), nil
	case Ref:
		return typeexpr.NewNamed(node.Reference()), nil
	case Primitive:
		return typeexpr.NewPrimitive(node.Kind()), nil
	case Or:
		return nil, errors.New(`expected a type, found "or"`)
	case Series:
		return nil, errors.New("expected a type, found a list separator")
	}
	return nil, errors.New("unknown type token")
}

// at reports whether the next token in cursor equals sep, without consuming it.
// Every separator token is an empty struct, so equality compares the variant.
func (p Parser) at(cursor *Cursor[Token], sep Token) bool {
	token, ok := cursor.Peek()
	return ok && token == sep
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import (
	"errors"

	"github.com/andreychh/tgen/model/types/v2"
)

// Parser decodes a stream of type tokens into a type expression by recursive
// descent. The grammar is:
//
//	expression := term (separator term)*
//	term       := "Array of" expression | ref | primitive
//
// Precedence is implicit in the rule layering: "Array of" recurses into a whole
// expression, so a union inside an array (the comma-and form) binds as the
// array's element, while a top-level union (the "or" form) joins terms.
type Parser struct {
	tokens *Cursor[Token]
}

// NewParser constructs a Parser over a cursor of type-expression tokens.
func NewParser(tokens *Cursor[Token]) Parser {
	return Parser{tokens: tokens}
}

// Expression returns the type expression decoded from the whole token stream. It
// fails when the tokens are malformed or do not all form a single expression.
//
//nolint:ireturn // types.Expression is the intentional decoded contract.
func (p Parser) Expression() (types.Expression, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.tokens.Done() {
		return nil, errors.New("trailing tokens after type expression")
	}
	return expr, nil
}

// expression parses a term, then folds any separator-joined terms that follow
// into a Union. It fails when a term is malformed.
//
//nolint:ireturn // types.Expression is the recursive node contract.
func (p Parser) expression() (types.Expression, error) {
	first, err := p.term()
	if err != nil {
		return nil, err
	}
	variants := []types.Expression{first}
	for p.atSeparator() {
		p.tokens.Skip()
		next, err := p.term()
		if err != nil {
			return nil, err
		}
		variants = append(variants, next)
	}
	if len(variants) == 1 {
		return variants[0], nil
	}
	return types.NewUnion(variants...), nil
}

// term parses an "Array of" wrapper around a nested expression, or a single
// reference or primitive. It fails on the end of input or a token that cannot
// begin a type.
//
//nolint:ireturn // types.Expression is the recursive node contract.
func (p Parser) term() (types.Expression, error) {
	node, ok := p.tokens.Take()
	if !ok {
		return nil, errors.New("expected a type, found end of input")
	}
	switch node := node.(type) {
	case ArrayOf:
		element, err := p.expression()
		if err != nil {
			return nil, err
		}
		return types.NewArray(element), nil
	case Ref:
		return types.NewNamed(node.Reference()), nil
	case Primitive:
		return types.NewPrimitive(node.Kind()), nil
	default:
		return nil, errors.New("expected a type")
	}
}

// atSeparator reports whether the next token is a separator, without consuming
// it.
func (p Parser) atSeparator() bool {
	token, ok := p.tokens.Peek()
	_, sep := token.(Separator)
	return ok && sep
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import (
	"errors"
	"fmt"
	"strings"
)

type TypeSource interface {
	AsString() (string, error)
}

// Type builds an expression tree from a field type source.
//
// The Telegram Bot API describes field types using human-readable English
// strings rather than a formal type syntax. Type parses these strings into
// a hierarchical AST that the rendering layer can traverse to generate target
// language types.
//
// Input conforms to the following grammar:
//
//	Type  ::= Union | Array | Named
//	Array ::= "Array" "of" Type
//	Union ::= Named ( "or" Named )+
//	            | Named ( "," Named )+ "and" Named
//	Named ::= <A-Za-z>+
//
// Union normalization: ", " and " and " are treated as equivalent to " or ".
// This covers all compound type forms used in the Telegram Bot API.
//
//nolint:dupword // "Named" is a grammar rule name, not a duplicate word
type Type struct {
	catalog Catalog
	source  TypeSource
}

func NewType(c Catalog, s TypeSource) Type {
	return Type{catalog: c, source: s}
}

func (t Type) AsExpression() (Expression, error) {
	value, err := t.source.AsString()
	if err != nil {
		return nil, err
	}
	return t.parse(value)
}

func (t Type) parse(expr string) (Expression, error) {
	if expr == "" {
		return nil, errors.New("unexpected empty type expression")
	}
	remainder, found := strings.CutPrefix(expr, "Array of ")
	if found {
		if remainder == "" {
			return nil, fmt.Errorf("incomplete array type: %q", expr)
		}
		inner, err := t.parse(remainder)
		if err != nil {
			return nil, err
		}
		return NewArray(inner), nil
	}
	normalized := strings.ReplaceAll(expr, " and ", " or ")
	normalized = strings.ReplaceAll(normalized, ", ", " or ")
	parts := strings.Split(normalized, " or ")
	if len(parts) > 1 {
		variants := make([]Expression, len(parts))
		for i, part := range parts {
			expr, err := t.parse(part)
			if err != nil {
				return nil, fmt.Errorf("parsing union variant %q: %w", part, err)
			}
			variants[i] = expr
		}
		return NewUnion(variants...), nil
	}
	kind, ok := t.catalog.Lookup(expr)
	if !ok {
		return nil, fmt.Errorf("unknown kind for %q", expr)
	}
	return NewNamed(expr, kind), nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"strings"
)

// TypeTree builds a TypeExpression tree from a FieldType.
//
// The Telegram Bot API describes field types using human-readable English
// strings rather than a formal type syntax. TypeTree parses these strings into
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
type TypeTree struct {
	source FieldType
}

// NewTypeTree creates a TypeTree from a FieldType.
func NewTypeTree(ft FieldType) TypeTree {
	return TypeTree{source: ft}
}

// Root parses the field type and returns the root of the type expression tree.
// Returns an error if the source is empty, contains invalid characters, or
// does not conform to the grammar.
//
//nolint:ireturn // TypeExpression is the intentional public contract of this method
func (t TypeTree) Root() (TypeExpression, error) {
	value, err := t.source.Value()
	if err != nil {
		return nil, err
	}
	return t.parse(value)
}

//nolint:ireturn // returns interface by design to match Root() contract
func (t TypeTree) parse(expr string) (TypeExpression, error) {
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
		return NewArrayType(inner), nil
	}
	normalized := strings.ReplaceAll(expr, " and ", " or ")
	normalized = strings.ReplaceAll(normalized, ", ", " or ")
	parts := strings.Split(normalized, " or ")
	if len(parts) > 1 {
		variants := make([]TypeExpression, len(parts))
		for i, part := range parts {
			expr, err := t.parse(part)
			if err != nil {
				return nil, fmt.Errorf("parsing union variant %q: %w", part, err)
			}
			variants[i] = expr
		}
		return NewUnionType(variants), nil
	}
	return NewNamedType(expr), nil
}

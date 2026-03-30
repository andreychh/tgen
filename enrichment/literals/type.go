// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package literals

import "github.com/andreychh/tgen/parsing/types"

// Type represents a parsing.Type wrapping a known TypeExpression.
type Type struct {
	expr types.TypeExpression
}

// NewType constructs a Type from expr.
func NewType(expr types.TypeExpression) Type {
	return Type{expr: expr}
}

func (t Type) AsExpression() (types.TypeExpression, error) {
	return t.expr, nil
}

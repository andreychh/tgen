// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package literals

import "github.com/andreychh/tgen/model/types"

// Type represents a model.Type wrapping a known TypeExpression.
type Type struct {
	expr types.Expression
}

// NewType constructs a Type from expr.
func NewType(expr types.Expression) Type {
	return Type{expr: expr}
}

func (t Type) AsExpression() (types.Expression, error) {
	return t.expr, nil
}

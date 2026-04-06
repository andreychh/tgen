// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/types"
)

//nolint:gochecknoglobals // immutable lookup table, not mutable global state
var primitives = map[string]string{
	"Integer": "int64",
	"Int":     "int64",
	"Float":   "float64",
	"String":  "string",
	"Boolean": "bool",
	"True":    "bool",
}

type ExprType struct {
	typ model.Type
}

func NewExprType(t model.Type) ExprType {
	return ExprType{typ: t}
}

func (t ExprType) AsString() (string, error) {
	expr, err := t.typ.AsExpression()
	if err != nil {
		return "", fmt.Errorf("getting type expr: %w", err)
	}
	return t.render(expr)
}

func (t ExprType) render(expr types.Expression) (string, error) {
	switch expr := expr.(type) {
	case types.Named:
		p, ok := primitives[expr.Name()]
		if !ok {
			return NewStringName(expr.Name()).AsString()
		}
		return p, nil
	case types.Array:
		elem, err := t.render(expr.Element())
		if err != nil {
			return "", err
		}
		return "[]" + elem, nil
	case types.Union:
		return "", fmt.Errorf("unsupported union %q", expr)
	}
	return "", fmt.Errorf("unknown type expression %q", expr)
}

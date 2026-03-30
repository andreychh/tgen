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

type Type struct {
	typ model.Type
}

func NewType(t model.Type) Type {
	return Type{typ: t}
}

func (t Type) AsString() (string, error) {
	expr, err := t.typ.AsExpression()
	if err != nil {
		return "", fmt.Errorf("getting type expr: %w", err)
	}
	return t.render(expr)
}

func (t Type) render(expr types.TypeExpression) (string, error) {
	if name, ok := expr.Named(); ok {
		p, ok := primitives[name]
		if !ok {
			return NewStringName(name).AsString()
		}
		return p, nil
	}
	if inner, ok := expr.Array(); ok {
		elem, err := t.render(inner)
		if err != nil {
			return "", err
		}
		return "[]" + elem, nil
	}
	if _, ok := expr.Union(); ok {
		return "any", nil
	}
	return "", fmt.Errorf("unknown type expression")
}

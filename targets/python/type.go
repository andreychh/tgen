// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/types"
)

//nolint:gochecknoglobals // immutable lookup table, not mutable global state
var parts = map[string]string{
	"Integer": "IntPart",
	"Int":     "IntPart",
	"Float":   "FloatPart",
	"String":  "StrPart",
	"Boolean": "BoolPart",
	"True":    "BoolPart",
}

//nolint:gochecknoglobals // immutable lookup table, not mutable global state
var primitives = map[string]string{
	"Integer": "int",
	"Int":     "int",
	"Float":   "float",
	"String":  "str",
	"Boolean": "bool",
	"True":    "bool",
}

type Type struct {
	inner model.Type
}

func NewType(t model.Type) Type {
	return Type{inner: t}
}

func (t Type) AsString() (string, error) {
	expr, err := t.inner.AsExpression()
	if err != nil {
		return "", fmt.Errorf("getting type expr: %w", err)
	}
	return t.render(expr)
}

func (t Type) Part() (string, error) {
	expr, err := t.inner.AsExpression()
	if err != nil {
		return "", fmt.Errorf("getting type expr: %w", err)
	}
	return t.part(expr)
}

func (t Type) part(expr types.Expression) (string, error) {
	switch e := expr.(type) {
	case types.Named:
		if p, ok := parts[e.Name()]; ok {
			return p + "(%s)", nil
		}
		return "%s", nil
	case types.Array:
		return "ListPart(%s)", nil
	case types.Union:
		return "", fmt.Errorf("unsupported union %q", expr)
	}
	return "", fmt.Errorf("unknown type expression %q", expr)
}

func (t Type) name() (string, error) {
	expr, err := t.inner.AsExpression()
	if err != nil {
		return "", fmt.Errorf("getting type expr: %w", err)
	}
	for {
		switch e := expr.(type) {
		case types.Named:
			return e.Name(), nil
		case types.Array:
			expr = e.Element()
		default:
			return "", fmt.Errorf("unknown type expression %q", expr)
		}
	}
}

func (t Type) render(expr types.Expression) (string, error) {
	switch expr := expr.(type) {
	case types.Named:
		p, ok := primitives[expr.Name()]
		if !ok {
			return NewStringClassName(expr.Name()).AsString()
		}
		return p, nil
	case types.Array:
		elem, err := t.render(expr.Element())
		if err != nil {
			return "", err
		}
		return "list[" + elem + "]", nil
	case types.Union:
		return "", fmt.Errorf("unsupported union %q", expr)
	}
	return "", fmt.Errorf("unknown type expression %q", expr)
}

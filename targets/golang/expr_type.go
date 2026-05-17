// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/model/types"
)

const (
	typeInteger = "Integer"
	typeInt     = "Int"
	typeFloat   = "Float"
	typeString  = "String"
	typeBoolean = "Boolean"
	typeTrue    = "True"
)

//nolint:gochecknoglobals // immutable lookup table, not mutable global state
var primitives = map[string]string{
	typeInteger: "int64",
	typeInt:     "int64",
	typeFloat:   "float64",
	typeString:  "string",
	typeBoolean: "bool",
	typeTrue:    "bool",
}

//nolint:gochecknoglobals // immutable lookup table, not mutable global state
var zeros = map[string]string{
	typeInteger: "0",
	typeInt:     "0",
	typeFloat:   "0",
	typeString:  `""`,
	typeBoolean: "false",
	typeTrue:    "false",
}

//nolint:gochecknoglobals // immutable lookup table, not mutable global state
var parts = map[string]string{
	typeInteger: "NewInt64Part",
	typeInt:     "NewInt64Part",
	typeFloat:   "NewFloat64Part",
	typeString:  "NewStringPart",
	typeBoolean: "NewBoolPart",
	typeTrue:    "NewBoolPart",
}

const zeroNil = "nil"

type ExprType struct {
	typ types.Expression
}

func NewExprType(t types.Expression) ExprType {
	return ExprType{typ: t}
}

func (t ExprType) IsPrimitive() (bool, error) {
	expr := t.typ
	for {
		switch e := expr.(type) {
		case types.Named:
			return e.Kind() == types.KindPrimitive, nil
		case types.Array:
			expr = e.Element()
		default:
			return false, fmt.Errorf("unexpected type expression %q", expr)
		}
	}
}

func (t ExprType) IsUnion() (bool, error) {
	expr := t.typ
	for {
		switch e := expr.(type) {
		case types.Named:
			return e.Kind() == types.KindUnion, nil
		case types.Array:
			expr = e.Element()
		default:
			return false, fmt.Errorf("unexpected type expression %q", expr)
		}
	}
}

func (t ExprType) Depth() (int, error) {
	expr := t.typ
	depth := 0
	for {
		switch e := expr.(type) {
		case types.Named:
			return depth, nil
		case types.Array:
			expr = e.Element()
			depth += 1
		default:
			return 0, fmt.Errorf("unexpected type expression %q", expr)
		}
	}
}

func (t ExprType) Name() (string, error) {
	expr := t.typ
	for {
		switch e := expr.(type) {
		case types.Named:
			return e.Name(), nil
		case types.Array:
			expr = e.Element()
		default:
			return "", fmt.Errorf("unexpected type expression %q", expr)
		}
	}
}

func (t ExprType) Part() (string, error) {
	depth, err := t.Depth()
	if err != nil {
		return "", err
	}
	name, err := t.Name()
	if err != nil {
		return "", err
	}
	if depth > 0 {
		if _, ok := parts[name]; ok {
			return "NewJSONPart(%s)", nil
		}
		return "NewSlicePart(%s)", nil
	}
	if part, ok := parts[name]; ok {
		return part + "(%s)", nil
	}
	return "%s", nil
}

func (t ExprType) Value() (string, error) {
	return t.render(t.typ)
}

func (t ExprType) Zero() (string, error) {
	depth, err := t.Depth()
	if err != nil {
		return "", err
	}
	if depth > 0 {
		return zeroNil, nil
	}
	isUnion, err := t.IsUnion()
	if err != nil {
		return "", err
	}
	if isUnion {
		return zeroNil, nil
	}
	name, err := t.Name()
	if err != nil {
		return "", err
	}
	if zero, ok := zeros[name]; ok {
		return zero, nil
	}
	formatted, err := NewStringName(name).Value()
	if err != nil {
		return "", err
	}
	return formatted + "{}", nil
}

func (t ExprType) render(expr types.Expression) (string, error) {
	switch expr := expr.(type) {
	case types.Named:
		p, ok := primitives[expr.Name()]
		if !ok {
			return NewStringName(expr.Name()).Value()
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

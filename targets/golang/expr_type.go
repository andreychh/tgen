// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"strings"

	"github.com/andreychh/tgen/model/ir"
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
	typ ir.Type
}

func NewExprType(typ ir.Type) ExprType {
	return ExprType{typ: typ}
}

func (t ExprType) IsPrimitive() (bool, error) {
	kind, err := t.typ.Kind()
	if err != nil {
		return false, err
	}
	return kind == types.KindPrimitive, nil
}

func (t ExprType) IsUnion() (bool, error) {
	kind, err := t.typ.Kind()
	if err != nil {
		return false, err
	}
	return kind == types.KindUnion, nil
}

func (t ExprType) Depth() (int, error) {
	return t.typ.Dimensionality()
}

func (t ExprType) Name() (string, error) {
	return t.typ.Name()
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
	name, err := t.typ.Name()
	if err != nil {
		return "", err
	}
	dim, err := t.typ.Dimensionality()
	if err != nil {
		return "", err
	}
	rendered, ok := primitives[name]
	if !ok {
		rendered, err = NewStringName(name).Value()
		if err != nil {
			return "", err
		}
	}
	return strings.Repeat("[]", dim) + rendered, nil
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

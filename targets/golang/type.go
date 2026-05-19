// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"strings"

	"github.com/andreychh/tgen/model"
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
	typeInteger: "newInt64Part",
	typeInt:     "newInt64Part",
	typeFloat:   "newFloat64Part",
	typeString:  "newStringPart",
	typeBoolean: "newBoolPart",
	typeTrue:    "newBoolPart",
}

const zeroNil = "nil"

// Type represents a Go type expression used in generated code, with its optionality.
type Type struct {
	typ ir.Type
	opt model.Optionality
}

// NewType creates a Type from an ir.Type and its optionality.
func NewType(typ ir.Type, opt model.Optionality) Type {
	return Type{typ: typ, opt: opt}
}

func NewRequiredType(typ ir.Type) Type {
	return NewType(typ, false)
}

func (t Type) IsUnion() (bool, error) {
	kind, err := t.typ.Kind()
	if err != nil {
		return false, err
	}
	return kind == types.KindUnion, nil
}

func (t Type) Shape() (Shape, error) {
	isUnion, err := t.IsUnion()
	if err != nil {
		return "", err
	}
	if !isUnion {
		return ShapePlain, nil
	}
	dim, err := t.typ.Dimensionality()
	if err != nil {
		return "", err
	}
	if dim == 1 {
		return ShapeUnionArray, nil
	}
	return ShapeUnion, nil
}

func (t Type) Name() (string, error) {
	return t.typ.Name()
}

func (t Type) Value() (string, error) {
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
		rendered = NewName(model.Name(name)).Value()
	}
	ptr, err := t.hasPointer()
	if err != nil {
		return "", err
	}
	if ptr {
		return "*" + rendered, nil
	}
	return strings.Repeat("[]", dim) + rendered, nil
}

func (t Type) Part() (string, error) {
	depth, err := t.typ.Dimensionality()
	if err != nil {
		return "", err
	}
	name, err := t.Name()
	if err != nil {
		return "", err
	}
	if depth > 0 {
		if _, ok := parts[name]; ok {
			return "newPrimitiveSlicePart(%s)", nil
		}
		return "newObjectSlicePart(%s)", nil
	}
	part, ok := parts[name]
	if !ok {
		return "%s", nil
	}
	ptr, err := t.hasPointer()
	if err != nil {
		return "", err
	}
	if ptr {
		return part + "(*%s)", nil
	}
	return part + "(%s)", nil
}

func (t Type) Zero() (string, error) {
	if t.opt {
		return zeroNil, nil
	}
	depth, err := t.typ.Dimensionality()
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
	formatted := NewName(model.Name(name)).Value()
	return formatted + "{}", nil
}

func (t Type) hasPointer() (bool, error) {
	if !t.opt {
		return false, nil
	}
	depth, err := t.typ.Dimensionality()
	if err != nil {
		return false, err
	}
	if depth > 0 {
		return false, nil
	}
	isUnion, err := t.IsUnion()
	if err != nil {
		return false, err
	}
	return !isUnion, nil
}

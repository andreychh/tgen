// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/ir"
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
	inner ir.Type
}

func NewType(t ir.Type) Type {
	return Type{inner: t}
}

func (t Type) Value() (string, error) {
	name, err := t.inner.Name()
	if err != nil {
		return "", err
	}
	dim, err := t.inner.Dimensionality()
	if err != nil {
		return "", err
	}
	rendered, ok := primitives[name]
	if !ok {
		rendered = NewClassName(model.Name(name)).Value()
	}
	for range dim {
		rendered = "list[" + rendered + "]"
	}
	return rendered, nil
}

func (t Type) Part() (string, error) {
	name, err := t.inner.Name()
	if err != nil {
		return "", err
	}
	dim, err := t.inner.Dimensionality()
	if err != nil {
		return "", err
	}
	part, isPrim := parts[name]
	if dim == 0 {
		if isPrim {
			return part + "(%s)", nil
		}
		return "%s", nil
	}
	if dim == 1 && !isPrim {
		return "ListFormJsonPart(%s)", nil
	}
	return "ListPart(%s)", nil
}

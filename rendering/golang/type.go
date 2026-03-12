// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"errors"
	"fmt"

	"github.com/andreychh/tgen/parsing"
)

//nolint:gochecknoglobals // immutable lookup table, not mutable global state
var namedTypes = map[string]string{
	"Integer": "int64",
	"Int":     "int64",
	"Float":   "float64",
	"String":  "string",
	"Boolean": "bool",
	"True":    "bool",
}

type Type struct {
	tree      parsing.TypeTree
	unionName RawValue
}

func NewType(t parsing.TypeTree, unionName RawValue) Type {
	return Type{tree: t, unionName: unionName}
}

func (t Type) Value() (string, error) {
	root, err := t.tree.Root()
	if err != nil {
		return "", fmt.Errorf("getting type root: %w", err)
	}
	return t.render(root)
}

func (t Type) render(expr parsing.TypeExpression) (string, error) {
	if name, ok := expr.Named(); ok {
		goName, ok := namedTypes[name]
		if !ok {
			return NewDefaultName(NewStaticName(name)).Value()
		}
		return goName, nil
	}
	if inner, ok := expr.Array(); ok {
		elem, err := t.render(inner)
		if err != nil {
			return "", err
		}
		return "[]" + elem, nil
	}
	if _, ok := expr.Union(); ok {
		return t.unionName.Value()
	}
	return "", errors.New("unknown type expression")
}

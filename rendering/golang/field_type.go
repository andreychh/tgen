// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"errors"
	"fmt"

	"github.com/andreychh/tgen/parsing"
)

var namedTypes = map[string]string{
	"Integer": "int64",
	"Float":   "float64",
	"String":  "string",
	"Boolean": "bool",
	"True":    "bool",
}

type FieldType struct {
	tree parsing.TypeTree
	key  parsing.FieldKey
}

func NewFieldType(t parsing.TypeTree, k parsing.FieldKey) FieldType {
	return FieldType{tree: t, key: k}
}

func (t FieldType) Value() (string, error) {
	root, err := t.tree.Root()
	if err != nil {
		return "", fmt.Errorf("getting type root: %w", err)
	}
	return t.render(root)
}

func (t FieldType) render(expr parsing.TypeExpression) (string, error) {
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
		return NewDefaultName(t.key).Value()
	}
	return "", errors.New("unknown type expression")
}

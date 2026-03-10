// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/parsing"
)

type OptionalFieldType struct {
	inner       FieldType
	tree        parsing.TypeTree
	optionality parsing.FieldOptionality
}

func NewOptionalFieldType(
	ft FieldType,
	t parsing.TypeTree,
	o parsing.FieldOptionality,
) OptionalFieldType {
	return OptionalFieldType{inner: ft, tree: t, optionality: o}
}

func (t OptionalFieldType) Value() (string, error) {
	typ, err := t.inner.Value()
	if err != nil {
		return "", fmt.Errorf("getting field type: %w", err)
	}
	optional, err := t.needsPointer()
	if err != nil {
		return "", err
	}
	if optional {
		return "*" + typ, nil
	}
	return typ, nil
}

func (t OptionalFieldType) needsPointer() (bool, error) {
	root, err := t.tree.Root()
	if err != nil {
		return false, err
	}
	_, isArray := root.Array()
	optional, err := t.optionality.Value()
	if err != nil {
		return false, fmt.Errorf("getting field optionality: %w", err)
	}
	return optional && !isArray, nil
}

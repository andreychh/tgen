// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/parsing"
)

type OptionalType struct {
	origin      Type
	tree        parsing.TypeTree
	optionality parsing.Optionality
}

func NewOptionalType(
	t Type,
	tree parsing.TypeTree,
	o parsing.Optionality,
) OptionalType {
	return OptionalType{origin: t, tree: tree, optionality: o}
}

func (t OptionalType) Value() (string, error) {
	typ, err := t.origin.Value()
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

func (t OptionalType) needsPointer() (bool, error) {
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

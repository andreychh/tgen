// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/model"
)

type OptionalType struct {
	typ         model.Type
	optionality model.Optionality
}

func NewOptionalType(t model.Type, o model.Optionality) OptionalType {
	return OptionalType{typ: t, optionality: o}
}

func (t OptionalType) AsString() (string, error) {
	typ, err := NewType(t.typ).AsString()
	if err != nil {
		return "", err
	}
	optional, err := t.optionality.AsBool()
	if err != nil {
		return "", fmt.Errorf("getting field optionality: %w", err)
	}
	if !optional {
		return typ, nil
	}
	expr, err := t.typ.AsExpression()
	if err != nil {
		return "", fmt.Errorf("getting type expr: %w", err)
	}
	_, isArray := expr.Array()
	if isArray {
		return typ, nil
	}
	return "*" + typ, nil
}

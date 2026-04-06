// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/types"
)

type OptionalType struct {
	typ         model.Type
	optionality model.Optionality
}

func NewOptionalType(t model.Type, o model.Optionality) OptionalType {
	return OptionalType{typ: t, optionality: o}
}

func (t OptionalType) AsString() (string, error) {
	typ, err := NewExprType(t.typ).AsString()
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
	if _, ok := expr.(types.Array); ok {
		return typ, nil
	}
	return "*" + typ, nil
}

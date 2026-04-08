// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/types"
)

type OptionalType struct {
	typ model.Type
	opt model.Optionality
}

func NewOptionalType(t model.Type, o model.Optionality) OptionalType {
	return OptionalType{typ: t, opt: o}
}

func (t OptionalType) IsUnion() (bool, error) {
	expr, err := t.typ.AsExpression()
	if err != nil {
		return false, fmt.Errorf("getting type expr: %w", err)
	}
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

func (t OptionalType) Depth() (int, error) {
	expr, err := t.typ.AsExpression()
	if err != nil {
		return 0, fmt.Errorf("getting type expr: %w", err)
	}
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

func (t OptionalType) Name() (string, error) {
	return NewExprType(t.typ).Name()
}

func (t OptionalType) AsString() (string, error) {
	typ, err := NewExprType(t.typ).AsString()
	if err != nil {
		return "", err
	}
	opt, err := t.opt.AsBool()
	if err != nil {
		return "", fmt.Errorf("getting field optionality: %w", err)
	}
	if !opt {
		return typ, nil
	}
	expr, err := t.typ.AsExpression()
	if err != nil {
		return "", fmt.Errorf("getting type expr: %w", err)
	}
	if _, ok := expr.(types.Array); ok {
		return typ, nil
	}
	isUnion, err := t.IsUnion()
	if err != nil {
		return "", err
	}
	if isUnion {
		return typ, nil
	}
	return "*" + typ, nil
}

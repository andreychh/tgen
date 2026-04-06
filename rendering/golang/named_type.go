// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/literals"
	"github.com/andreychh/tgen/model/types"
)

type NamedType struct {
	name model.Name
}

func NewNamedType(name model.Name) NamedType {
	return NamedType{name: name}
}

func (t NamedType) AsString() (string, error) {
	name, err := t.name.AsString()
	if err != nil {
		return "", fmt.Errorf("getting name: %w", err)
	}
	return NewExprType(literals.NewType(types.NewNamed(name, types.KindUnknown))).AsString()
}

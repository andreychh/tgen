// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package flattened

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/corrected"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeform"
)

// Alias is a name tgen introduces for a type the documentation leaves unnamed.
// Its type is reduced to a flat one.
type Alias struct {
	Ref         model.Reference
	Name        model.Name
	Type        typeform.Type
	Description prose.Passage
}

// AliasMapping maps a corrected alias into a flattened one by reducing its type
// expression to a flat type.
type AliasMapping struct{}

// NewAliasMapping constructs an AliasMapping.
func NewAliasMapping() AliasMapping {
	return AliasMapping{}
}

// Apply implements [pipeline.Mapping]. It fails when the alias's type holds a
// union.
func (m AliasMapping) Apply(alias corrected.Alias) (Alias, error) {
	typ, err := NewNarrowing(alias.Type).Value()
	if err != nil {
		return Alias{}, fmt.Errorf("narrowing alias type: %w", err)
	}
	return Alias{
		Ref:         alias.Ref,
		Name:        alias.Name,
		Type:        typ,
		Description: alias.Description,
	}, nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typebound"
)

// Alias is the record of a name tgen gives a type the documentation leaves
// unnamed, joined with the type that name stands for.
type Alias struct {
	Ref         model.Reference
	Name        model.Name
	Type        typebound.Type
	Description prose.Passage
}

func (Alias) isDefinition() {}

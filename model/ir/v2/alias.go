// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typebound"
)

// Alias is the record of a name tgen gives a type the documentation leaves
// unnamed, joined with the type that name stands for. Direction is which way
// the alias travels between a client and the API; the declaration it stands for
// is the same either way, since an alias is another name for a type that
// already answers for its own encoding.
type Alias struct {
	Ref         model.Reference
	Name        model.Name
	Type        typebound.Type
	Description prose.Passage
	Direction   model.Direction
}

func (Alias) isDefinition() {}

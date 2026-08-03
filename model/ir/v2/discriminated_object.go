// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
)

// DiscriminatedObject is the record of an object a union tells apart by a
// discriminator, joined with the fields it owns and the fixed value it is told
// apart by. It stands beside [Object], not on top of it: the two partition the
// objects the documentation names, and no object is both. Files narrows Fields
// to those reaching a file, and is empty for an object holding none. Rewrites
// reports the same as it does for [Object]. Direction is which way the object
// travels between a client and the API: one no request carries is never
// encoded, so the fixed value telling it apart never has to be written at all,
// and one no response carries is never decoded. Introduced reports that tgen
// introduced the object rather than reading it from the documentation page,
// which leaves it no section a target can address.
type DiscriminatedObject struct {
	Ref           model.Reference
	Name          model.Name
	Description   prose.Passage
	Fields        []Field
	Files         []FileField
	Rewrites      bool
	Direction     model.Direction
	Introduced    bool
	Discriminator Discriminator
}

func (DiscriminatedObject) isDefinition() {}

// Discriminator is the field an object is told apart by and the fixed value it
// holds there. The pipeline moves that field out of the field table, so this
// value is the only place it survives.
type Discriminator struct {
	Key   model.Key
	Value model.DiscriminatorValue
}

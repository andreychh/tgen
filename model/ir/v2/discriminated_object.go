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
// reports the same as it does for [Object].
type DiscriminatedObject struct {
	Ref           model.Reference
	Name          model.Name
	Description   prose.Passage
	Fields        []Field
	Files         []FileField
	Rewrites      bool
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

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

// DiscriminatedObject represents a variant of a discriminated union parsed from
// the specification.
type DiscriminatedObject struct {
	root, h4 gq.Selection
}

// NewDiscriminatedObject constructs a DiscriminatedObject from a h4 selection.
func NewDiscriminatedObject(root, h4 gq.Selection) DiscriminatedObject {
	return DiscriminatedObject{root: root, h4: h4}
}

func (o DiscriminatedObject) Name() model.Name {
	return NewName(o.h4)
}

func (o DiscriminatedObject) Fields() explicit.Fields {
	return NewDiscriminatedObjectFields(o.root, o.h4)
}

func (o DiscriminatedObject) Reference() model.Reference {
	return NewDefinitionReference(o.h4.Find("a.anchor"))
}

func (o DiscriminatedObject) Description() model.Description {
	return NewDefinitionDescription(o.h4)
}

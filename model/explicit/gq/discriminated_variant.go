// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

// DiscriminatedVariant represents a variant of a discriminated union parsed from the
// specification.
type DiscriminatedVariant struct {
	h4 gq.Selection
}

// NewDiscriminatedVariant constructs a DiscriminatedVariant from a h4 selection.
func NewDiscriminatedVariant(h4 gq.Selection) DiscriminatedVariant {
	return DiscriminatedVariant{h4: h4}
}

func (v DiscriminatedVariant) Name() model.Name {
	return NewName(v.h4)
}

func (v DiscriminatedVariant) Fields() explicit.Fields {
	return NewFields(v.h4)
}

func (v DiscriminatedVariant) Reference() model.Reference {
	return NewDefinitionReference(v.h4.Find("a.anchor"))
}

func (v DiscriminatedVariant) Description() model.Description {
	return NewDefinitionDescription(v.h4)
}

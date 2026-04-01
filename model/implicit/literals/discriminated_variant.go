// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package literals

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/literals"
)

type DiscriminatedVariant struct {
	name  string
	value string
}

// NewDiscriminatedVariant constructs a DiscriminatedVariant from a field name.
func NewDiscriminatedVariant(name, value string) DiscriminatedVariant {
	return DiscriminatedVariant{name: name, value: value}
}

func (v DiscriminatedVariant) Name() model.Name {
	return literals.NewName(v.name)
}

func (v DiscriminatedVariant) DiscriminatorValue() model.DiscriminatorValue {
	return literals.NewDiscriminatorValue(v.value)
}

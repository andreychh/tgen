// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"slices"

	"github.com/andreychh/tgen/model/types"
	"github.com/andreychh/tgen/pkg/gq"
)

type Catalog struct {
	root gq.Selection
}

func NewCatalog(root gq.Selection) Catalog {
	return Catalog{root: root}
}

func (c Catalog) Lookup(name string) (types.Kind, bool) {
	if c.isPrimitive(name) {
		return types.KindPrimitive, true
	}
	definition := c.root.Find("div#dev_page_content h4").
		FilterFunc(func(h4 gq.Selection) bool {
			return h4.Text() == name
		}).
		At(0)
	if definition.IsEmpty() {
		return types.KindUnknown, false
	}
	switch NewHeader(c.root, definition).Kind() {
	case DefinitionKindObject,
		DefinitionKindDiscriminatedVariant:
		return types.KindObject, true
	case DefinitionKindDiscriminatedUnion,
		DefinitionKindStructuredUnion,
		DefinitionKindFallbackUnion:
		return types.KindUnion, true
	case DefinitionKindMethod,
		DefinitionKindUnknown:
		return types.KindUnknown, false
	}
	return types.KindUnknown, false
}

func (c Catalog) isPrimitive(name string) bool {
	return slices.Contains(
		[]string{"Integer", "Int", "Float", "String", "Boolean", "True"},
		name,
	)
}

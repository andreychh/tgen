// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

// GQDiscriminatedUnion represents a union whose variants are distinguished by a
// fixed-value discriminator field (e.g. type="emoji"). It navigates from the
// document root to locate each variant's full definition.
type GQDiscriminatedUnion struct {
	root      gq.Selection
	selection gq.Selection
}

// NewDiscriminatedUnion constructs a GQDiscriminatedUnion from the document root
// and the h4 selection of the union definition.
func NewDiscriminatedUnion(root, h4 gq.Selection) GQDiscriminatedUnion {
	return GQDiscriminatedUnion{root: root, selection: h4}
}

func (u GQDiscriminatedUnion) Ref() DefinitionRef {
	return NewDefinitionRef(u.selection.Find("a.anchor"))
}

func (u GQDiscriminatedUnion) Name() ObjectName {
	return NewGQObjectName(u.selection)
}

func (u GQDiscriminatedUnion) Description() DefinitionDescription {
	return NewDefinitionDescription(u.selection)
}

// Variants returns the variant objects of this union. Each variant name is read
// from the union's <ul> list and resolved to a full object definition by navigating
// from the document root.
func (u GQDiscriminatedUnion) Variants() iter.Seq[VariantObject] {
	return func(yield func(VariantObject) bool) {
		for li := range u.selection.Until("h3, h4, hr").Find("ul li").All() {
			if !yield(NewVariantObject(u.root.
				Find("div#dev_page_content h4").
				FilterFunc(func(s gq.Selection) bool {
					return s.Text() == li.Find("a").Text()
				}).
				At(0),
			)) {
				break
			}
		}
	}
}

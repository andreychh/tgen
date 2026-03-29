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
	root gq.Selection
	h4   gq.Selection
}

// NewGQDiscriminatedUnion constructs a GQDiscriminatedUnion from the document
// root and the h4 td of the union definition.
func NewGQDiscriminatedUnion(root, h4 gq.Selection) GQDiscriminatedUnion {
	return GQDiscriminatedUnion{root: root, h4: h4}
}

func (u GQDiscriminatedUnion) Reference() Reference {
	return NewGQDefinitionReference(u.h4.Find("a.anchor"))
}

func (u GQDiscriminatedUnion) Name() Name {
	return NewGQName(u.h4)
}

func (u GQDiscriminatedUnion) Description() Description {
	return NewGQDefinitionDescription(u.h4)
}

// DiscriminatorKey returns the field key shared by all variants of this union
// (e.g. "type" for ReactionType).
//
//nolint:varnamelen // <h4> is the standard HTML heading element name
func (u GQDiscriminatedUnion) DiscriminatorKey() Key {
	li := u.h4.Until("h3, h4, hr").Find("ul li").At(0)
	h4 := u.root.
		Find("div#dev_page_content h4").
		FilterFunc(func(s gq.Selection) bool {
			return s.Text() == li.Find("a").Text()
		}).
		At(0)
	return NewGQDiscriminatedVariant(h4).Fields().Discriminator().Key()
}

// Variants returns the variant objects of this union. Each variant name is read
// from the union's <ul> list and resolved to a full object definition by
// navigating from the document root.
//
//nolint:varnamelen // <h4> is the standard HTML heading element name
func (u GQDiscriminatedUnion) Variants() iter.Seq[DiscriminatedVariant] {
	return func(yield func(DiscriminatedVariant) bool) {
		for li := range u.h4.Until("h3, h4, hr").Find("ul li").All() {
			h4 := u.root.
				Find("div#dev_page_content h4").
				FilterFunc(func(s gq.Selection) bool {
					return s.Text() == li.Find("a").Text()
				}).
				At(0)
			if !yield(NewGQDiscriminatedVariant(h4)) {
				break
			}
		}
	}
}

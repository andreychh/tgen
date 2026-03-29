// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

type GQUnions struct {
	root gq.Selection
}

func NewGQUnions(root gq.Selection) GQUnions {
	return GQUnions{root: root}
}

func (u GQUnions) Discriminated() iter.Seq[DiscriminatedUnion] {
	return func(yield func(DiscriminatedUnion) bool) {
		seq := u.root.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewGQHeader(u.root, h4).Kind() == DefinitionKindDiscriminatedUnion
			}).
			All()
		for h4 := range seq {
			if !yield(NewGQDiscriminatedUnion(u.root, h4)) {
				break
			}
		}
	}
}

func (u GQUnions) Structured() iter.Seq[StructuredUnion] {
	return func(yield func(StructuredUnion) bool) {
		panic("not implemented")
	}
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

type Unions struct {
	root gq.Selection
}

func NewUnions(root gq.Selection) Unions {
	return Unions{root: root}
}

func (u Unions) Discriminated() iter.Seq[explicit.DiscriminatedUnion] {
	return func(yield func(explicit.DiscriminatedUnion) bool) {
		seq := u.root.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewHeader(u.root, h4).Kind() == DefinitionKindDiscriminatedUnion
			}).
			All()
		for h4 := range seq {
			if !yield(NewDiscriminatedUnion(u.root, h4)) {
				break
			}
		}
	}
}

func (u Unions) Structured() iter.Seq[explicit.StructuredUnion] {
	return func(yield func(explicit.StructuredUnion) bool) {
		seq := u.root.
			Find("div#dev_page_content h4").
			FilterFunc(func(h4 gq.Selection) bool {
				return NewHeader(u.root, h4).Kind() == DefinitionKindStructuredUnion
			}).
			All()
		for h4 := range seq {
			if !yield(NewStructuredUnion(u.root, h4)) {
				break
			}
		}
	}
}

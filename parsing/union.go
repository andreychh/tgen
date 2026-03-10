// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

type Union struct {
	selection gq.Selection
}

func NewUnion(h4 gq.Selection) Union {
	return Union{selection: h4}
}

func (u Union) Ref() DefinitionRef {
	return NewDefinitionRef(u.selection.Find("a.anchor"))
}

func (u Union) Name() ObjectName {
	return NewObjectName(u.selection)
}

func (u Union) Description() DefinitionDescription {
	return NewDefinitionDescription(u.selection)
}

func (u Union) Variants() iter.Seq[Variant] {
	return func(yield func(Variant) bool) {
		seq := u.selection.
			Until("h3, h4, hr").
			Find("ul li").
			All()
		for li := range seq {
			if !yield(NewVariant(li)) {
				break
			}
		}
	}
}

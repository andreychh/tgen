// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

type GQStructuredUnion struct {
	h4 gq.Selection
}

func NewGQStructuredUnion(h4 gq.Selection) GQStructuredUnion {
	return GQStructuredUnion{h4: h4}
}

func (u GQStructuredUnion) Reference() Reference {
	return NewGQDefinitionReference(u.h4.Find("a.anchor"))
}

func (u GQStructuredUnion) Name() Name {
	return NewGQName(u.h4)
}

func (u GQStructuredUnion) Description() Description {
	return NewGQDefinitionDescription(u.h4)
}

func (u GQStructuredUnion) Variants() iter.Seq[StructuredVariant] {
	return func(yield func(StructuredVariant) bool) {
		seq := u.h4.
			Until("h3, h4, hr").
			Find("ul li").
			All()
		for li := range seq {
			if !yield(NewGQStructuredVariant(li)) {
				break
			}
		}
	}
}

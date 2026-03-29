// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

type GQObject struct {
	h4 gq.Selection
}

func NewGQObject(h4 gq.Selection) GQObject {
	return GQObject{h4: h4}
}

func (o GQObject) Reference() Reference {
	return NewGQDefinitionReference(o.h4.Find("a.anchor"))
}

func (o GQObject) Name() Name {
	return NewGQName(o.h4)
}

func (o GQObject) Description() Description {
	return NewGQDefinitionDescription(o.h4)
}

func (o GQObject) Fields() iter.Seq[Field] {
	return func(yield func(Field) bool) {
		seq := o.h4.
			Until("h3, h4, hr").
			Find("table tbody tr").
			All()
		for tr := range seq {
			if !yield(NewGQObjectField(tr)) {
				break
			}
		}
	}
}

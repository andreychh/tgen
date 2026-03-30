// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"iter"

	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

type Object struct {
	h4 gq.Selection
}

func NewObject(h4 gq.Selection) Object {
	return Object{h4: h4}
}

func (o Object) Reference() explicit.Reference {
	return NewDefinitionReference(o.h4.Find("a.anchor"))
}

func (o Object) Name() explicit.Name {
	return NewName(o.h4)
}

func (o Object) Description() explicit.Description {
	return NewDefinitionDescription(o.h4)
}

func (o Object) Fields() iter.Seq[explicit.Field] {
	return func(yield func(explicit.Field) bool) {
		seq := o.h4.
			Until("h3, h4, hr").
			Find("table tbody tr").
			All()
		for tr := range seq {
			if !yield(NewObjectField(tr)) {
				break
			}
		}
	}
}

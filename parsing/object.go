// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

type Object struct {
	selection gq.Selection
}

func NewObject(h4 gq.Selection) Object {
	return Object{selection: h4}
}

func (o Object) Ref() DefinitionRef {
	return NewDefinitionRef(o.selection.Find("a.anchor"))
}

//nolint:ireturn // ObjectName is the intentional public contract of this method
func (o Object) Name() ObjectName {
	return NewGQObjectName(o.selection)
}

func (o Object) Description() GQDefinitionDescription {
	return NewDefinitionDescription(o.selection)
}

func (o Object) Fields() iter.Seq[ObjectField] {
	return func(yield func(ObjectField) bool) {
		seq := o.selection.
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

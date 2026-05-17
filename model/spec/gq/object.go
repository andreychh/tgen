// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/pkg/gq"
)

type Object struct {
	root, h4 gq.Selection
}

func NewObject(root, h4 gq.Selection) Object {
	return Object{root: root, h4: h4}
}

func (o Object) Reference() (model.Reference, error) {
	return NewDefinitionReference(o.h4.Find("a.anchor")).Value()
}

func (o Object) Name() (model.Name, error) {
	return NewName(o.h4).Value()
}

func (o Object) Description() model.Description {
	return NewDefinitionDescription(o.h4)
}

func (o Object) Fields() iter.Seq[spec.Field] {
	return func(yield func(spec.Field) bool) {
		seq := o.h4.
			Until("h3, h4, hr").
			Find("table tbody tr").
			All()
		for tr := range seq {
			if !yield(NewObjectField(o.root, tr)) {
				break
			}
		}
	}
}

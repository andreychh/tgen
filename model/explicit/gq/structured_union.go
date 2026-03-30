// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq //nolint:dupl // structurally identical to Object but semantically distinct: StructuredUnion has list variants, Object has table fields

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

type StructuredUnion struct {
	h4 gq.Selection
}

func NewStructuredUnion(h4 gq.Selection) StructuredUnion {
	return StructuredUnion{h4: h4}
}

func (u StructuredUnion) Reference() model.Reference {
	return NewDefinitionReference(u.h4.Find("a.anchor"))
}

func (u StructuredUnion) Name() model.Name {
	return NewName(u.h4)
}

func (u StructuredUnion) Description() model.Description {
	return NewDefinitionDescription(u.h4)
}

func (u StructuredUnion) Variants() iter.Seq[explicit.StructuredVariant] {
	return func(yield func(explicit.StructuredVariant) bool) {
		seq := u.h4.
			Until("h3, h4, hr").
			Find("ul li").
			All()
		for li := range seq {
			if !yield(NewStructuredVariant(li)) {
				break
			}
		}
	}
}

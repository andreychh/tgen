// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

// Method represents a Telegram Bot API method definition parsed from HTML.
type Method struct {
	root, h4 gq.Selection
}

// NewMethod creates a Method from an h4 selection.
func NewMethod(root, h4 gq.Selection) Method {
	return Method{root: root, h4: h4}
}

func (m Method) Reference() model.Reference {
	return NewDefinitionReference(m.h4.Find("a.anchor"))
}

func (m Method) Name() model.Name {
	return NewName(m.h4)
}

func (m Method) Description() model.Description {
	return NewDefinitionDescription(m.h4)
}

func (m Method) ReturnType() model.Type {
	return NewReturnType(m.root, m.h4)
}

func (m Method) Fields() iter.Seq[explicit.Field] {
	return func(yield func(explicit.Field) bool) {
		seq := m.h4.
			Until("h3, h4, hr").
			Find("table tbody tr").
			All()
		for tr := range seq {
			if !yield(NewMethodField(m.root, tr)) {
				break
			}
		}
	}
}

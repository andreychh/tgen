// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

// GQMethod represents a Telegram Bot API method definition parsed from HTML.
type GQMethod struct {
	h4 gq.Selection
}

// NewGQMethod creates a GQMethod from an h4 selection.
func NewGQMethod(h4 gq.Selection) GQMethod {
	return GQMethod{h4: h4}
}

func (m GQMethod) Reference() Reference {
	return NewGQDefinitionReference(m.h4.Find("a.anchor"))
}

func (m GQMethod) Name() Name {
	return NewGQName(m.h4)
}

func (m GQMethod) Description() Description {
	return NewGQDefinitionDescription(m.h4)
}

func (m GQMethod) ReturnType() Type {
	return NewGQReturnType(m.h4)
}

func (m GQMethod) Fields() iter.Seq[Field] {
	return func(yield func(Field) bool) {
		seq := m.h4.
			Until("h3, h4, hr").
			Find("table tbody tr").
			All()
		for tr := range seq {
			if !yield(NewGQMethodField(tr)) {
				break
			}
		}
	}
}

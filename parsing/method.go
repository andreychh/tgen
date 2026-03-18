// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

// GQMethod represents a Telegram Bot API method definition parsed from HTML.
type GQMethod struct {
	selection gq.Selection
}

// NewMethod creates a GQMethod from an h4 selection.
func NewMethod(h4 gq.Selection) GQMethod {
	return GQMethod{selection: h4}
}

func (m GQMethod) Ref() DefinitionRef {
	return NewDefinitionRef(m.selection.Find("a.anchor"))
}

func (m GQMethod) Name() MethodName {
	return NewMethodName(m.selection)
}

func (m GQMethod) Description() GQDefinitionDescription {
	return NewDefinitionDescription(m.selection)
}

//nolint:ireturn // TypeTree is the intentional public contract of this method
func (m GQMethod) Returns() TypeTree {
	return NewReturnType(m.selection)
}

func (m GQMethod) Fields() iter.Seq[MethodField] {
	return func(yield func(MethodField) bool) {
		seq := m.selection.
			Until("h3, h4, hr").
			Find("table tbody tr").
			All()
		for tr := range seq {
			if !yield(NewMethodField(tr)) {
				break
			}
		}
	}
}

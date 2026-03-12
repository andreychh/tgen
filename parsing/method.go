// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"

	"github.com/andreychh/tgen/parsing/gq"
)

type Method struct {
	selection gq.Selection
}

func NewMethod(h4 gq.Selection) Method {
	return Method{selection: h4}
}

func (m Method) Ref() DefinitionRef {
	return NewDefinitionRef(m.selection.Find("a.anchor"))
}

func (m Method) Name() MethodName {
	return NewMethodName(m.selection)
}

func (m Method) Description() DefinitionDescription {
	return NewDefinitionDescription(m.selection)
}

func (m Method) Returns() TypeTree {
	return NewTypeTree(NewReturnType(m.selection))
}

func (m Method) Fields() iter.Seq[MethodField] {
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

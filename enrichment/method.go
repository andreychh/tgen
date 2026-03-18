// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
)

// Method represents an enriched Telegram Bot API method definition.
type Method struct {
	inner parsing.Method
}

// NewMethod constructs a Method from a parsed method with all method rules applied.
func NewMethod(m parsing.Method) Method {
	return Method{inner: NewMaybeMessageRule(m)}
}

func (m Method) Ref() parsing.DefinitionRef {
	return m.inner.Ref()
}

func (m Method) Name() parsing.MethodName {
	return m.inner.Name()
}

//nolint:ireturn // DefinitionDescription is the intentional public contract of this method
func (m Method) Description() parsing.DefinitionDescription {
	return m.inner.Description()
}

//nolint:ireturn // TypeTree is the intentional public contract of this method
func (m Method) Returns() parsing.TypeTree {
	return m.inner.Returns()
}

func (m Method) Fields() iter.Seq[parsing.MethodField] {
	return m.inner.Fields()
}

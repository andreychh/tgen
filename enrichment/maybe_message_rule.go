// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
)

// MaybeMessageRule represents an enrichment rule that replaces Message-or-True return types with MaybeMessage.
type MaybeMessageRule struct {
	inner parsing.Method
}

func NewMaybeMessageRule(m parsing.Method) MaybeMessageRule {
	return MaybeMessageRule{inner: m}
}

func (r MaybeMessageRule) Ref() parsing.DefinitionRef {
	return r.inner.Ref()
}

func (r MaybeMessageRule) Name() parsing.MethodName {
	return r.inner.Name()
}

func (r MaybeMessageRule) Description() parsing.GQDefinitionDescription {
	return r.inner.Description()
}

func (r MaybeMessageRule) Fields() iter.Seq[parsing.MethodField] {
	return r.inner.Fields()
}

//nolint:ireturn // TypeTree is the intentional public contract of this method
func (r MaybeMessageRule) Returns() parsing.TypeTree {
	root, err := r.inner.Returns().Root()
	if err != nil {
		return r.inner.Returns()
	}
	if !root.Equal(parsing.NewUnionType([]parsing.TypeExpression{
		parsing.NewNamedType("Message"),
		parsing.NewNamedType("True"),
	})) {
		return r.inner.Returns()
	}
	return parsing.NewTypeTreeExpr(parsing.NewNamedType("MaybeMessage"))
}

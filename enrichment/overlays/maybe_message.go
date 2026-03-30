// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/enrichment/literals"
	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/parsing/types"
)

// MaybeMessage represents a parsing.Method overlay that replaces
// Message-or-True return types with MaybeMessage.
type MaybeMessage struct {
	inner parsing.Method
}

// NewMaybeMessage constructs a MaybeMessage from m.
func NewMaybeMessage(m parsing.Method) MaybeMessage {
	return MaybeMessage{inner: m}
}

func (r MaybeMessage) Reference() parsing.Reference     { return r.inner.Reference() }
func (r MaybeMessage) Name() parsing.Name               { return r.inner.Name() }
func (r MaybeMessage) Description() parsing.Description { return r.inner.Description() }
func (r MaybeMessage) Fields() iter.Seq[parsing.Field]  { return r.inner.Fields() }

func (r MaybeMessage) ReturnType() parsing.Type {
	expr, err := r.inner.ReturnType().AsExpression()
	if err != nil {
		return r.inner.ReturnType()
	}
	if !expr.Equals(types.NewUnionType([]types.TypeExpression{
		types.NewNamedType("Message"),
		types.NewNamedType("True"),
	})) {
		return r.inner.ReturnType()
	}
	return literals.NewType(types.NewNamedType("MaybeMessage"))
}

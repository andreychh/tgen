// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/parsing/explicit/literals"
	"github.com/andreychh/tgen/parsing/types"
)

// MaybeMessage represents a explicit.Method overlay that replaces
// Message-or-True return types with MaybeMessage.
type MaybeMessage struct {
	inner explicit.Method
}

// NewMaybeMessage constructs a MaybeMessage from m.
func NewMaybeMessage(m explicit.Method) MaybeMessage {
	return MaybeMessage{inner: m}
}

func (r MaybeMessage) Reference() explicit.Reference     { return r.inner.Reference() }
func (r MaybeMessage) Name() explicit.Name               { return r.inner.Name() }
func (r MaybeMessage) Description() explicit.Description { return r.inner.Description() }
func (r MaybeMessage) Fields() iter.Seq[explicit.Field]  { return r.inner.Fields() }

func (r MaybeMessage) ReturnType() explicit.Type {
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

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/model/literals"
	"github.com/andreychh/tgen/model/types"
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

func (r MaybeMessage) Reference() model.Reference {
	return r.inner.Reference()
}

func (r MaybeMessage) Name() model.Name {
	return r.inner.Name()
}

func (r MaybeMessage) Description() model.Description {
	return r.inner.Description()
}

func (r MaybeMessage) ReturnType() model.Type {
	expr, err := r.inner.ReturnType().AsExpression()
	if err != nil {
		return r.inner.ReturnType()
	}
	if !expr.Equals(types.NewUnion(
		types.NewNamed("Message", types.KindObject),
		types.NewNamed("True", types.KindPrimitive),
	)) {
		return r.inner.ReturnType()
	}
	return literals.NewType(types.NewNamed("MaybeMessage", types.KindUnion))
}

func (r MaybeMessage) Fields() iter.Seq[explicit.Field] {
	return r.inner.Fields()
}

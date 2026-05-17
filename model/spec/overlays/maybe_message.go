// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/model/types"
)

// MaybeMessage represents a spec.Method overlay that replaces
// Message-or-True return types with MaybeMessage.
type MaybeMessage struct {
	inner spec.Method
}

// NewMaybeMessage constructs a MaybeMessage from m.
func NewMaybeMessage(m spec.Method) MaybeMessage {
	return MaybeMessage{inner: m}
}

func (r MaybeMessage) Reference() (model.Reference, error) {
	return r.inner.Reference()
}

func (r MaybeMessage) Name() (model.Name, error) {
	return r.inner.Name()
}

func (r MaybeMessage) Description() model.Description {
	return r.inner.Description()
}

func (r MaybeMessage) ReturnType() (types.Expression, error) {
	expr, err := r.inner.ReturnType()
	if err != nil {
		return nil, err
	}
	if !expr.Equals(types.NewUnion(
		types.NewNamed("Message", types.KindObject),
		types.NewNamed("True", types.KindPrimitive),
	)) {
		return expr, nil
	}
	return types.NewNamed("MaybeMessage", types.KindUnion), nil
}

func (r MaybeMessage) Fields() iter.Seq[spec.Field] {
	return r.inner.Fields()
}

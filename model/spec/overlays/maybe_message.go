// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"fmt"
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

func (r MaybeMessage) Result() (spec.Result, error) {
	result, err := r.inner.Result()
	if err != nil {
		return nil, err
	}
	switch result := result.(type) {
	case spec.Command:
		return result, nil
	case spec.Value:
		if !result.Type().Equals(types.NewUnion(
			types.NewNamed("Message", types.KindObject),
			types.NewNamed("True", types.KindPrimitive),
		)) {
			return result, nil
		}
		return spec.NewValue(types.NewNamed("MaybeMessage", types.KindUnion)), nil
	default:
		return nil, fmt.Errorf("classifying method result: unexpected result type %T", result)
	}
}

func (r MaybeMessage) Fields() iter.Seq[spec.Field] {
	return r.inner.Fields()
}

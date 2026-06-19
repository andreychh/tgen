// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"fmt"
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/pkg/iters"
)

// Method represents a Telegram Bot API method definition narrowed for code
// generation.
type Method struct {
	inner spec.Method
}

// NewMethod constructs a Method from a parsed method.
func NewMethod(m spec.Method) Method {
	return Method{inner: m}
}

func (m Method) Reference() (model.Reference, error) {
	return m.inner.Reference()
}

func (m Method) Name() (model.Name, error) {
	return m.inner.Name()
}

func (m Method) Description() model.Description {
	return m.inner.Description()
}

func (m Method) Result() (Result, error) {
	result, err := m.inner.Result()
	if err != nil {
		return nil, err
	}
	switch result := result.(type) {
	case spec.Command:
		return NewCommand(), nil
	case spec.Value:
		return NewValue(NewType(result.Type())), nil
	default:
		return nil, fmt.Errorf("narrowing method result: unexpected result type %T", result)
	}
}

func (m Method) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(m.inner.Fields(), NewField)
}

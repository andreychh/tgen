// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
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

func (m Method) ReturnType() (Type, error) {
	expr, err := m.inner.ReturnType()
	if err != nil {
		return Type{}, err
	}
	return NewType(expr), nil
}

func (m Method) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(m.inner.Fields(), NewField)
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

type Method struct {
	inner explicit.Method
}

func NewMethod(o explicit.Method) Method {
	return Method{inner: o}
}

func (m Method) Name() Name {
	return NewName(m.inner.Name())
}

func (m Method) Doc() GoDoc {
	return NewGoDoc(NewDefinitionDoc(m.inner.Reference(), m.inner.Description()))
}

func (m Method) ReturnType() Type {
	return NewExprType(m.inner.ReturnType())
}

func (m Method) Fields() iter.Seq[Field] {
	return iters.NewMappedSeq(m.inner.Fields(), NewField)
}

func (m Method) IsMultipart() bool {
	return iters.IsAny(m.inner.Fields(), func(f explicit.Field) bool {
		name, err := NewExprType(f.Type()).Name()
		if err != nil {
			return false
		}
		return name == "InputFile"
	})
}

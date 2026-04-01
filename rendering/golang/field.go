// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/model/explicit"

type Field struct {
	inner explicit.Field
}

func NewField(f explicit.Field) Field {
	return Field{inner: f}
}

func (f Field) Name() Name {
	return NewName(f.inner.Key())
}

func (f Field) Type() Type {
	return NewOptionalType(f.inner.Type(), f.inner.Optionality())
}

func (f Field) Tag() Tag {
	return NewTag(f.inner.Key(), f.inner.Optionality())
}

func (f Field) Doc() GoDoc {
	return NewGoDoc(f.inner.Description())
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/parsing"

type Field struct {
	inner parsing.Field
}

func NewField(f parsing.Field) Field {
	return Field{inner: f}
}

func (f Field) Name() Name {
	return NewDefaultName(f.inner.Key())
}

func (f Field) Type() OptionalFieldType {
	return NewOptionalFieldType(
		NewFieldType(f.inner.Type(), f.inner.Key()),
		f.inner.Type(),
		f.inner.IsOptional(),
	)
}

func (f Field) Tag() FieldTag {
	return NewFieldTag(f.inner.Key(), f.inner.IsOptional())
}

func (f Field) Doc() Doc {
	return NewDoc(f.inner.Description())
}

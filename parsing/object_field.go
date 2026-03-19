// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"github.com/andreychh/tgen/parsing/gq"
)

type ObjectField struct {
	selection gq.Selection
}

func NewObjectField(tr gq.Selection) ObjectField {
	return ObjectField{selection: tr}
}

func (f ObjectField) Key() FieldKey {
	return NewFieldKey(f.selection.Find("td").At(0))
}

//nolint:ireturn // TypeTree is the intentional public contract of Field
func (f ObjectField) Type() TypeTree {
	return NewTypeTree(NewFieldType(f.selection.Find("td").At(1)))
}

//nolint:ireturn // Optionality is the intentional public contract of Field
func (f ObjectField) IsOptional() Optionality {
	return NewObjectFieldOptionality(f.selection.Find("td").At(2))
}

//nolint:ireturn // FieldDescription is the intentional public contract of Field
func (f ObjectField) Description() FieldDescription {
	return NewObjectFieldDescription(f.selection.Find("td").At(2))
}

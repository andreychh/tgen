// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import "github.com/andreychh/tgen/parsing/gq"

type MethodField struct {
	selection gq.Selection
}

func NewMethodField(tr gq.Selection) MethodField {
	return MethodField{selection: tr}
}

func (f MethodField) Key() FieldKey {
	return NewFieldKey(f.selection.Find("td").At(0))
}

//nolint:ireturn // TypeTree is the intentional public contract of Field
func (f MethodField) Type() TypeTree {
	return NewTypeTree(NewFieldType(f.selection.Find("td").At(1)))
}

//nolint:ireturn // Optionality is the intentional public contract of Field
func (f MethodField) IsOptional() Optionality {
	return NewMethodFieldOptionality(f.selection.Find("td").At(2))
}

//nolint:ireturn // FieldDescription is the intentional public contract of Field
func (f MethodField) Description() FieldDescription {
	return NewMethodFieldDescription(f.selection.Find("td").At(3))
}

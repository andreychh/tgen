// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"github.com/andreychh/tgen/parsing/gq"
	"github.com/andreychh/tgen/parsing/types"
)

type GQObjectField struct {
	tr gq.Selection
}

func NewGQObjectField(tr gq.Selection) GQObjectField {
	return GQObjectField{tr: tr}
}

func (f GQObjectField) Key() Key {
	return NewGQKey(f.tr.Find("td").At(0))
}

func (f GQObjectField) Type() Type {
	return types.NewType(NewGQType(f.tr.Find("td").At(1)))
}

func (f GQObjectField) Optionality() Optionality {
	return NewGQObjectFieldOptionality(f.tr.Find("td").At(2))
}

func (f GQObjectField) Description() Description {
	return NewGQObjectFieldDescription(f.tr.Find("td").At(2))
}

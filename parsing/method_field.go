// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"github.com/andreychh/tgen/parsing/gq"
	"github.com/andreychh/tgen/parsing/types"
)

type GQMethodField struct {
	tr gq.Selection
}

func NewGQMethodField(tr gq.Selection) GQMethodField {
	return GQMethodField{tr: tr}
}

func (f GQMethodField) Key() Key {
	return NewGQKey(f.tr.Find("td").At(0))
}

func (f GQMethodField) Type() Type {
	return types.NewType(NewGQType(f.tr.Find("td").At(1)))
}

func (f GQMethodField) Optionality() Optionality {
	return NewGQMethodFieldOptionality(f.tr.Find("td").At(2))
}

func (f GQMethodField) Description() Description {
	return NewGQMethodFieldDescription(f.tr.Find("td").At(3))
}

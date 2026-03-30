// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/parsing/types"
	"github.com/andreychh/tgen/pkg/gq"
)

type ObjectField struct {
	tr gq.Selection
}

func NewObjectField(tr gq.Selection) ObjectField {
	return ObjectField{tr: tr}
}

func (f ObjectField) Key() explicit.Key {
	return NewKey(f.tr.Find("td").At(0))
}

func (f ObjectField) Type() explicit.Type {
	return types.NewType(NewType(f.tr.Find("td").At(1)))
}

func (f ObjectField) Optionality() explicit.Optionality {
	return NewObjectFieldOptionality(f.tr.Find("td").At(2))
}

func (f ObjectField) Description() explicit.Description {
	return NewObjectFieldDescription(f.tr.Find("td").At(2))
}

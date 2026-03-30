// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/types"
	"github.com/andreychh/tgen/pkg/gq"
)

type MethodField struct {
	tr gq.Selection
}

func NewMethodField(tr gq.Selection) MethodField {
	return MethodField{tr: tr}
}

func (f MethodField) Key() model.Key {
	return NewKey(f.tr.Find("td").At(0))
}

func (f MethodField) Type() model.Type {
	return types.NewType(NewType(f.tr.Find("td").At(1)))
}

func (f MethodField) Optionality() model.Optionality {
	return NewMethodFieldOptionality(f.tr.Find("td").At(2))
}

func (f MethodField) Description() model.Description {
	return NewMethodFieldDescription(f.tr.Find("td").At(3))
}

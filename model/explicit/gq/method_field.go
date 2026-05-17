// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/types"
	"github.com/andreychh/tgen/pkg/gq"
)

type MethodField struct {
	root, tr gq.Selection
}

func NewMethodField(root, tr gq.Selection) MethodField {
	return MethodField{root: root, tr: tr}
}

func (f MethodField) Key() (model.Key, error) {
	return NewKey(f.tr.Find("td").At(0)).Value()
}

func (f MethodField) Type() model.Type {
	return types.NewType(NewCatalog(f.root), NewType(f.tr.Find("td").At(1)))
}

func (f MethodField) Optionality() (model.Optionality, error) {
	return NewMethodFieldOptionality(f.tr.Find("td").At(2)).Value()
}

func (f MethodField) Description() model.Description {
	return NewMethodFieldDescription(f.tr.Find("td").At(3))
}

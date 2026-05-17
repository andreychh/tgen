// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/types"
	"github.com/andreychh/tgen/pkg/gq"
)

type ObjectField struct {
	root, tr gq.Selection
}

func NewObjectField(root, tr gq.Selection) ObjectField {
	return ObjectField{root: root, tr: tr}
}

func (f ObjectField) Key() (model.Key, error) {
	return NewKey(f.tr.Find("td").At(0)).Value()
}

func (f ObjectField) Type() model.Type {
	return types.NewType(NewCatalog(f.root), NewType(f.tr.Find("td").At(1)))
}

func (f ObjectField) Optionality() (model.Optionality, error) {
	return NewObjectFieldOptionality(f.tr.Find("td").At(2)).Value()
}

func (f ObjectField) Description() model.Description {
	return NewObjectFieldDescription(f.tr.Find("td").At(2))
}

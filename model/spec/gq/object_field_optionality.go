// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"strings"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/pkg/gq"
)

type ObjectFieldOptionality struct {
	td gq.Selection
}

func NewObjectFieldOptionality(td gq.Selection) ObjectFieldOptionality {
	return ObjectFieldOptionality{td: td}
}

// Value returns whether the object field is optional.
func (o ObjectFieldOptionality) Value() (model.Optionality, error) {
	if o.td.IsEmpty() {
		return false, errors.New("description column not found")
	}
	return model.Optionality(strings.HasPrefix(o.td.Text(), "Optional. ")), nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/pkg/gq"
)

type MethodFieldOptionality struct {
	td gq.Selection
}

func NewMethodFieldOptionality(td gq.Selection) MethodFieldOptionality {
	return MethodFieldOptionality{td: td}
}

// Value returns whether the method field is optional.
func (o MethodFieldOptionality) Value() (model.Optionality, error) {
	if o.td.IsEmpty() {
		return false, errors.New("required column not found")
	}
	return model.Optionality(o.td.Text() == "Optional"), nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"

	"github.com/andreychh/tgen/pkg/gq"
)

type MethodFieldOptionality struct {
	td gq.Selection
}

func NewMethodFieldOptionality(td gq.Selection) MethodFieldOptionality {
	return MethodFieldOptionality{td: td}
}

func (o MethodFieldOptionality) AsBool() (bool, error) {
	if o.td.IsEmpty() {
		return false, errors.New("required column not found")
	}
	return o.td.Text() == "Optional", nil
}

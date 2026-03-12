// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"

	"github.com/andreychh/tgen/parsing/gq"
)

type MethodFieldOptionality struct {
	selection gq.Selection
}

func NewMethodFieldOptionality(td gq.Selection) MethodFieldOptionality {
	return MethodFieldOptionality{selection: td}
}

func (o MethodFieldOptionality) Value() (bool, error) {
	if o.selection.IsEmpty() {
		return false, errors.New("required column not found")
	}
	return o.selection.Text() == "Optional", nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"

	"github.com/andreychh/tgen/parsing/gq"
)

type GQMethodFieldOptionality struct {
	td gq.Selection
}

func NewGQMethodFieldOptionality(td gq.Selection) GQMethodFieldOptionality {
	return GQMethodFieldOptionality{td: td}
}

func (o GQMethodFieldOptionality) AsBool() (bool, error) {
	if o.td.IsEmpty() {
		return false, errors.New("required column not found")
	}
	return o.td.Text() == "Optional", nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"strings"

	"github.com/andreychh/tgen/parsing/gq"
)

type GQObjectFieldOptionality struct {
	td gq.Selection
}

func NewGQObjectFieldOptionality(td gq.Selection) GQObjectFieldOptionality {
	return GQObjectFieldOptionality{td: td}
}

func (o GQObjectFieldOptionality) AsBool() (bool, error) {
	if o.td.IsEmpty() {
		return false, errors.New("description column not found")
	}
	return strings.HasPrefix(o.td.Text(), "Optional. "), nil
}

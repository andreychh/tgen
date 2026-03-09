// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"strings"

	"github.com/andreychh/tgen/parsing/gq"
)

type FieldOptionality struct {
	selection gq.Selection
}

func NewFieldOptionality(td gq.Selection) FieldOptionality {
	return FieldOptionality{selection: td}
}

func (d FieldOptionality) Value() (bool, error) {
	if d.selection.IsEmpty() {
		return false, errors.New("description column not found")
	}
	return strings.HasPrefix(d.selection.Text(), "Optional. "), nil
}

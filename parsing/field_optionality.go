// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"strings"

	"github.com/andreychh/tgen/parsing/gq"
)

type ObjectFieldOptionality struct {
	selection gq.Selection
}

func NewObjectFieldOptionality(td gq.Selection) ObjectFieldOptionality {
	return ObjectFieldOptionality{selection: td}
}

func (o ObjectFieldOptionality) Value() (bool, error) {
	if o.selection.IsEmpty() {
		return false, errors.New("description column not found")
	}
	return strings.HasPrefix(o.selection.Text(), "Optional. "), nil
}

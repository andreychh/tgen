// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"

	"github.com/andreychh/tgen/parsing/gq"
)

type MethodFieldDescription struct {
	selection gq.Selection
}

func NewMethodFieldDescription(td gq.Selection) MethodFieldDescription {
	return MethodFieldDescription{selection: td}
}

func (d MethodFieldDescription) Value() (string, error) {
	if d.selection.IsEmpty() {
		return "", errors.New("description column not found")
	}
	return d.selection.Text(), nil
}

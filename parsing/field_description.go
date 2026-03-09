// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"strings"

	"github.com/andreychh/tgen/parsing/gq"
)

type FieldDescription struct {
	selection gq.Selection
}

func NewFieldDescription(td gq.Selection) FieldDescription {
	return FieldDescription{selection: td}
}

func (d FieldDescription) Value() (string, error) {
	if d.selection.IsEmpty() {
		return "", errors.New("description column not found")
	}
	return strings.TrimPrefix(d.selection.Text(), "Optional. "), nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"strings"

	"github.com/andreychh/tgen/parsing/gq"
)

type ObjectFieldDescription struct {
	selection gq.Selection
}

func NewObjectFieldDescription(td gq.Selection) ObjectFieldDescription {
	return ObjectFieldDescription{selection: td}
}

func (d ObjectFieldDescription) Value() (string, error) {
	if d.selection.IsEmpty() {
		return "", errors.New("description column not found")
	}
	return strings.TrimPrefix(d.selection.Text(), "Optional. "), nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var fieldTypeRegex = regexp.MustCompile(`^[a-zA-Z0-9 ,]+$`)

type FieldType struct {
	selection gq.Selection
}

func NewFieldType(td gq.Selection) FieldType {
	return FieldType{selection: td}
}

func (f FieldType) Value() (string, error) {
	val := f.selection.Text()
	if !fieldTypeRegex.MatchString(val) {
		return "", fmt.Errorf("invalid field type: %q", val)
	}
	return val, nil
}

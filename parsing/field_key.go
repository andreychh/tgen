// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var fieldKeyRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

type FieldKey struct {
	selection gq.Selection
}

func NewFieldKey(td gq.Selection) FieldKey {
	return FieldKey{selection: td}
}

func (k FieldKey) Value() (string, error) {
	val := k.selection.Text()
	if !fieldKeyRegex.MatchString(val) {
		return "", fmt.Errorf("invalid field key: %q", val)
	}
	return val, nil
}

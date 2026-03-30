// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var typeRegex = regexp.MustCompile(`^[a-zA-Z0-9 ,]+$`)

type GQType struct {
	td gq.Selection
}

func NewGQType(td gq.Selection) GQType {
	return GQType{td: td}
}

func (t GQType) AsString() (string, error) {
	if t.td.IsEmpty() {
		return "", errors.New("field type not found")
	}
	typ := t.td.Text()
	if !typeRegex.MatchString(typ) {
		return "", fmt.Errorf("field type %q contains invalid characters", typ)
	}
	return typ, nil
}

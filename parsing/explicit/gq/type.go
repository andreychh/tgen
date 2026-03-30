// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/pkg/gq"
)

var typeRegex = regexp.MustCompile(`^[a-zA-Z0-9 ,]+$`)

type Type struct {
	td gq.Selection
}

func NewType(td gq.Selection) Type {
	return Type{td: td}
}

func (t Type) AsString() (string, error) {
	if t.td.IsEmpty() {
		return "", errors.New("field type not found")
	}
	typ := t.td.Text()
	if !typeRegex.MatchString(typ) {
		return "", fmt.Errorf("field type %q contains invalid characters", typ)
	}
	return typ, nil
}

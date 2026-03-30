// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"strings"

	"github.com/andreychh/tgen/pkg/gq"
)

type ObjectFieldOptionality struct {
	td gq.Selection
}

func NewObjectFieldOptionality(td gq.Selection) ObjectFieldOptionality {
	return ObjectFieldOptionality{td: td}
}

func (o ObjectFieldOptionality) AsBool() (bool, error) {
	if o.td.IsEmpty() {
		return false, errors.New("description column not found")
	}
	return strings.HasPrefix(o.td.Text(), "Optional. "), nil
}

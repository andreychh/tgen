// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"github.com/andreychh/tgen/model"
)

//nolint:gochecknoglobals // immutable lookup table, not mutable global state
var keywords = map[string]bool{
	"from": true, "import": true, "class": true, "return": true,
}

type FieldName struct {
	key model.Key
}

func NewFieldName(k model.Key) FieldName {
	return FieldName{key: k}
}

func (n FieldName) Value() string {
	name := string(n.key)
	if keywords[name] {
		return name + "_"
	}
	return name
}

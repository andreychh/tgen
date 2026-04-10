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
	inner model.Name
}

func NewFieldName(n model.Name) FieldName {
	return FieldName{inner: n}
}

func (n FieldName) AsString() (string, error) {
	name, err := n.inner.AsString()
	if err != nil {
		return "", err
	}
	if keywords[name] {
		return name + "_", nil
	}
	return name, nil
}

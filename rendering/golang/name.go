// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"strings"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/literals"
	"github.com/iancoleman/strcase"
)

//nolint:gochecknoglobals // immutable lookup table, not mutable global state
var acronyms = map[string]string{
	"Id":  "ID",
	"Url": "URL",
	"Api": "API",
	"Ip":  "IP",
}

type Name struct {
	inner model.Name
}

func NewName(n model.Name) Name {
	return Name{inner: n}
}

func NewStringName(s string) Name {
	return NewName(literals.NewName(s))
}

func (n Name) AsString() (string, error) {
	val, err := n.inner.AsString()
	if err != nil {
		return "", err
	}
	camel := strcase.ToCamel(val)
	for wrong, right := range acronyms {
		camel = strings.ReplaceAll(camel, wrong, right)
	}
	return camel, nil
}

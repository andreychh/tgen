// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

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

type ClassName struct {
	inner model.Name
}

func NewClassName(n model.Name) ClassName {
	return ClassName{inner: n}
}

func NewStringClassName(s string) ClassName {
	return NewClassName(literals.NewName(s))
}

func (n ClassName) AsString() (string, error) {
	name, err := n.inner.AsString()
	if err != nil {
		return "", err
	}
	camel := strcase.ToCamel(name)
	for wrong, right := range acronyms {
		camel = strings.ReplaceAll(camel, wrong, right)
	}
	return camel, nil
}

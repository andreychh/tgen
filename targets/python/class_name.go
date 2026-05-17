// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"strings"

	"github.com/andreychh/tgen/model"
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
	return NewClassName(model.Name(s))
}

func (n ClassName) Value() (string, error) {
	camel := strcase.ToCamel(string(n.inner))
	for wrong, right := range acronyms {
		camel = strings.ReplaceAll(camel, wrong, right)
	}
	return camel, nil
}

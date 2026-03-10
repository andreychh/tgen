// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"strings"

	"github.com/iancoleman/strcase"
)

//nolint:gochecknoglobals // immutable lookup table, not mutable global state
var defaultAcronyms = map[string]string{
	"Id":  "ID",
	"Url": "URL",
	"Api": "API",
	"Ip":  "IP",
}

type Name struct {
	inner    RawValue
	acronyms map[string]string
}

func NewName(n RawValue, acronyms map[string]string) Name {
	return Name{
		inner:    n,
		acronyms: acronyms,
	}
}

func NewDefaultName(n RawValue) Name {
	return NewName(n, defaultAcronyms)
}

func (n Name) Value() (string, error) {
	val, err := n.inner.Value()
	if err != nil {
		return "", err
	}
	camel := strcase.ToCamel(val)
	for wrong, right := range n.acronyms {
		camel = strings.ReplaceAll(camel, wrong, right)
	}
	return camel, nil
}

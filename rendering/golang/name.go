// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"strings"

	"github.com/iancoleman/strcase"
)

var defaultAcronyms = map[string]string{
	"Id":  "ID",
	"Url": "URL",
	"Api": "API",
	"Ip":  "IP",
}

type RawName interface {
	Value() (string, error)
}

type Name struct {
	inner    RawName
	acronyms map[string]string
}

func NewName(n RawName, acronyms map[string]string) Name {
	return Name{
		inner:    n,
		acronyms: acronyms,
	}
}

func NewDefaultName(n RawName) Name {
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

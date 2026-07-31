// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golangv2

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

// Name represents a documentation name rendered as a Go identifier.
type Name struct {
	inner model.Name
}

// NewName creates a Name from a documentation name.
func NewName(n model.Name) Name {
	return Name{inner: n}
}

// NewNameFromKey creates a Name from a field key.
func NewNameFromKey(k model.Key) Name {
	return NewName(model.Name(k))
}

// Value returns the name in Go's exported camel case, with the acronyms Go
// spells in capitals restored.
func (n Name) Value() string {
	camel := strcase.ToCamel(string(n.inner))
	for wrong, right := range acronyms {
		camel = strings.ReplaceAll(camel, wrong, right)
	}
	return camel
}

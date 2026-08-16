// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

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

// Name represents a documentation name rendered as a Python class name.
type Name struct {
	inner model.Name
}

// NewName creates a Name from a documentation name.
func NewName(n model.Name) Name {
	return Name{inner: n}
}

// Value returns the name in the capitalized words a class is declared by, with
// the acronyms spelled in capitals. The Go target spells them the same way, and
// a name is what a reader of either target looks the other one up by.
func (n Name) Value() string {
	camel := strcase.ToCamel(string(n.inner))
	for wrong, right := range acronyms {
		camel = strings.ReplaceAll(camel, wrong, right)
	}
	return camel
}

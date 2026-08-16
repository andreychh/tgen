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

// ClassName represents a documentation name rendered as the name of a Python
// class. It stands apart from [FieldName] because Python spells the two
// differently, which is what the Go target is spared: Go writes an exported
// name one way wherever it stands, and renders both through one view.
type ClassName struct {
	inner model.Name
}

// NewClassName creates a ClassName from a documentation name.
func NewClassName(n model.Name) ClassName {
	return ClassName{inner: n}
}

// Value returns the name in the capitalized words a class is declared by, with
// the acronyms spelled in capitals. The Go target spells them the same way, and
// a name is what a reader of either target looks the other one up by.
func (n ClassName) Value() string {
	camel := strcase.ToCamel(string(n.inner))
	for wrong, right := range acronyms {
		camel = strings.ReplaceAll(camel, wrong, right)
	}
	return camel
}

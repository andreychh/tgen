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
// the acronyms spelled in capitals, and followed by an underscore where the
// result is a word Python reserves. The Go target spells the capitals the same
// way, and a name is what a reader of either target looks the other one up by.
//
// Only True, False and None can ever need the underscore: every other reserved
// word is lowercase, and none survives the capitalization as itself. tgen names
// one definition True today, and Python would read the declaration of a class
// by that name as an assignment to a literal.
func (n ClassName) Value() string {
	camel := strcase.ToCamel(string(n.inner))
	for wrong, right := range acronyms {
		camel = strings.ReplaceAll(camel, wrong, right)
	}
	if keywords[model.Key(camel)] {
		return camel + "_"
	}
	return camel
}

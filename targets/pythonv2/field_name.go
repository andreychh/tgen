// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"github.com/andreychh/tgen/model"
)

// keywords holds every word Python reserves, which is what a field key may not
// be spelled as. It holds the reserved words alone: a soft keyword — match,
// case, type, _ — is a name the grammar reads as a keyword only where a keyword
// can stand, and remains a legal field name everywhere else. Spelling one of
// those away would rename "type", the key half the documented objects are told
// apart by.
//
//nolint:gochecknoglobals // immutable lookup table, not mutable global state
var keywords = map[model.Key]bool{
	"False": true, "None": true, "True": true, "and": true, "as": true,
	"assert": true, "async": true, "await": true, "break": true, "class": true,
	"continue": true, "def": true, "del": true, "elif": true, "else": true,
	"except": true, "finally": true, "for": true, "from": true, "global": true,
	"if": true, "import": true, "in": true, "is": true, "lambda": true,
	"nonlocal": true, "not": true, "or": true, "pass": true, "raise": true,
	"return": true, "try": true, "while": true, "with": true, "yield": true,
}

// FieldName represents a field key rendered as the name of a Python field. It
// stands apart from [ClassName] because Python spells the two differently: a
// class is declared by capitalized words, a field by the lowercase ones the
// documentation already gives.
type FieldName struct {
	key model.Key
}

// NewFieldName creates a FieldName from a field key.
func NewFieldName(k model.Key) FieldName {
	return FieldName{key: k}
}

// Value returns the key as the documentation writes it, followed by an
// underscore where the key is a word Python reserves. Nothing else is done to
// it: a key is already spelled in the lowercase words Python declares a field
// by. The key itself survives as the alias [Field.Assignment] writes.
func (n FieldName) Value() string {
	if keywords[n.key] {
		return string(n.key) + "_"
	}
	return string(n.key)
}

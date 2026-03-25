// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var definitionRowDiscriminatorRegex = regexp.MustCompile("always \u201c[^\u201d]+\u201d|must be [a-z][a-z0-9_]*\\s*$")

// FieldRowKind identifies the role of a field table row in a variant object.
type FieldRowKind string

const (
	KindFreeField          FieldRowKind = "free"
	KindDiscriminatorField FieldRowKind = "discriminator"
)

// DefinitionRow classifies a table row in a variant object definition.
type DefinitionRow struct {
	selection gq.Selection
}

// NewDefinitionRow constructs a DefinitionRow from a table row selection.
func NewDefinitionRow(tr gq.Selection) DefinitionRow {
	return DefinitionRow{selection: tr}
}

// Kind returns KindDiscriminatorField if the row's description contains an
// always-quoted value, and KindFreeField otherwise.
func (r DefinitionRow) Kind() FieldRowKind {
	if definitionRowDiscriminatorRegex.MatchString(r.selection.Find("td").At(2).Text()) {
		return KindDiscriminatorField
	}
	return KindFreeField
}

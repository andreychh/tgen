// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var fieldRowDiscriminatorRegex = regexp.MustCompile(
	"always \u201c[^\u201d]+\u201d|must be [a-z][a-z0-9_]*\\s*$",
)

// FieldKind identifies the role of a field table row in a variant object.
type FieldKind string

const (
	FieldKindFree          FieldKind = "free"
	FieldKindDiscriminator FieldKind = "discriminator"
)

// GQFieldRow classifies a table row in a variant object definition.
type GQFieldRow struct {
	tr gq.Selection
}

// NewGQFieldRow constructs a GQFieldRow from a table row td.
func NewGQFieldRow(tr gq.Selection) GQFieldRow {
	return GQFieldRow{tr: tr}
}

// Kind returns FieldKindDiscriminator if the row's description contains an
// always-quoted value, and FieldKindFree otherwise.
func (r GQFieldRow) Kind() FieldKind {
	if fieldRowDiscriminatorRegex.MatchString(r.tr.Find("td").At(2).Text()) {
		return FieldKindDiscriminator
	}
	return FieldKindFree
}

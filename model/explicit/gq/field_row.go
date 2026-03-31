// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"regexp"

	"github.com/andreychh/tgen/pkg/gq"
)

var fieldRowDiscriminatorRegex = regexp.MustCompile(
	"[Aa]lways (?:\u201c[^\u201d]+\u201d|\\d+\\.)|must be [a-z][a-z0-9_]*\\s*$",
)

// FieldKind identifies the role of a field table row in a variant object.
type FieldKind string

const (
	FieldKindFree          FieldKind = "free"
	FieldKindDiscriminator FieldKind = "discriminator"
)

// FieldRow classifies a table row in a variant object definition.
type FieldRow struct {
	tr gq.Selection
}

// NewFieldRow constructs a FieldRow from a table row td.
func NewFieldRow(tr gq.Selection) FieldRow {
	return FieldRow{tr: tr}
}

// Kind returns FieldKindDiscriminator if the row's description contains an
// always-quoted value, and FieldKindFree otherwise.
func (r FieldRow) Kind() FieldKind {
	if fieldRowDiscriminatorRegex.MatchString(r.tr.Find("td").At(2).Text()) {
		return FieldKindDiscriminator
	}
	return FieldKindFree
}

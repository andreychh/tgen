// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var discriminatorValueRegex = regexp.MustCompile(
	"always \u201c([^\u201d]+)\u201d|must be ([a-z][a-z0-9_]*)\\s*$",
)

// GQDiscriminator represents the discriminator field of a variant, parsed from
// its table row.
type GQDiscriminator struct {
	selection gq.Selection
}

// NewDiscriminator constructs a GQDiscriminator from a table row selection.
func NewDiscriminator(tr gq.Selection) GQDiscriminator {
	return GQDiscriminator{selection: tr}
}

// Key returns the field key of the discriminator field.
func (d GQDiscriminator) Key() FieldKey {
	return NewFieldKey(d.selection.Find("td").At(0))
}

// Value returns the fixed discriminator value from the field description.
// It returns an error if the description contains no always-quoted value.
func (d GQDiscriminator) Value() (string, error) {
	match := discriminatorValueRegex.FindStringSubmatch(d.selection.Find("td").At(2).Text())
	if match == nil {
		return "", errors.New("no discriminator value in description")
	}
	if match[1] != "" {
		return match[1], nil
	}
	return match[2], nil
}

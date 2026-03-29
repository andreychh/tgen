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

// GQDiscriminatorValue represents the fixed value of a discriminator field,
// parsed from the field description cell.
type GQDiscriminatorValue struct {
	td gq.Selection
}

// NewGQDiscriminatorValue constructs a GQDiscriminatorValue from a description td.
func NewGQDiscriminatorValue(td gq.Selection) GQDiscriminatorValue {
	return GQDiscriminatorValue{td: td}
}

// AsString returns the discriminator value extracted from the description.
// Returns an error if no always-quoted or must-be value is found.
func (v GQDiscriminatorValue) AsString() (string, error) {
	match := discriminatorValueRegex.FindStringSubmatch(v.td.Text())
	if match == nil {
		return "", errors.New("no discriminator value in description")
	}
	if match[1] != "" {
		return match[1], nil
	}
	return match[2], nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"regexp"

	"github.com/andreychh/tgen/pkg/gq"
)

var discriminatorValueRegex = regexp.MustCompile(
	"[Aa]lways (?:\u201c([^\u201d]+)\u201d|(\\d+)\\.)|must be ([a-z][a-z0-9_]*)\\s*$",
)

// DiscriminatorValue represents the fixed value of a discriminator field,
// parsed from the field description cell.
type DiscriminatorValue struct {
	td gq.Selection
}

// NewDiscriminatorValue constructs a DiscriminatorValue from a description td.
func NewDiscriminatorValue(td gq.Selection) DiscriminatorValue {
	return DiscriminatorValue{td: td}
}

// AsString returns the discriminator value extracted from the description.
// Returns an error if no always-quoted, always-numeric, or must-be value is found.
func (v DiscriminatorValue) AsString() (string, error) {
	match := discriminatorValueRegex.FindStringSubmatch(v.td.Text())
	if match == nil {
		return "", errors.New("no discriminator value in description")
	}
	for _, m := range match[1:] {
		if m != "" {
			return m, nil
		}
	}
	return "", errors.New("no discriminator value in description")
}

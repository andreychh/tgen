// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/parsing/dom"
)

// RawVariant implements Variant by wrapping an <li> element.
//
// It expects the selection to point to a list item containing an anchor link to
// the object definition.
type RawVariant struct {
	selection dom.Selection
	idRegex   *regexp.Regexp
	nameRegex *regexp.Regexp
}

// NewRawVariant creates a new RawVariant instance with custom validation
// patterns.
func NewRawVariant(s dom.Selection, id, name *regexp.Regexp) RawVariant {
	return RawVariant{
		selection: s,
		idRegex:   id,
		nameRegex: name,
	}
}

// NewDefaultRawVariant creates a new RawVariant instance using the default
// Telegram Bot API validation patterns.
func NewDefaultRawVariant(s dom.Selection) RawVariant {
	return NewRawVariant(s, idRegex, nameRegex)
}

// ID returns the anchor href found in the <li> element.
//
// It returns an error if the anchor tag is missing or if the extracted value
// does not match the configured ID pattern.
func (v RawVariant) ID() (string, error) {
	val, exists := v.selection.Find("a").Attr("href")
	if !exists {
		return "", errors.New("attribute href not found")
	}
	if !v.idRegex.MatchString(val) {
		return "", fmt.Errorf("id %q does not match pattern %s", val, v.idRegex)
	}
	return val, nil
}

// Name returns the text content of the <li> element.
//
// It returns an error if the content does not match the configured Name
// pattern.
func (v RawVariant) Name() (string, error) {
	val := v.selection.Text()
	if !v.nameRegex.MatchString(val) {
		return "", fmt.Errorf("name %q does not match pattern %s", val, v.nameRegex)
	}
	return val, nil
}

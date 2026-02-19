// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"iter"
	"regexp"
	"strings"

	"github.com/andreychh/tgen/parsing/dom"
)

// RawUnion implements Union by wrapping an <h4> element.
//
// It expects the selection to point to a definition header, immediately
// followed by description paragraphs and a <ul> list of variants.
type RawUnion struct {
	selection dom.Selection
	idRegex   *regexp.Regexp
	nameRegex *regexp.Regexp
}

// NewRawUnion creates a new RawUnion instance with custom validation patterns.
func NewRawUnion(s dom.Selection, id, name *regexp.Regexp) RawUnion {
	return RawUnion{
		selection: s,
		idRegex:   id,
		nameRegex: name,
	}
}

// NewDefaultRawUnion creates a new RawUnion instance using the default Telegram
// Bot API validation patterns.
func NewDefaultRawUnion(s dom.Selection) RawUnion {
	return NewRawUnion(s, idRegex, nameRegex)
}

// ID returns the anchor href found in the <h4> element.
//
// It returns an error if the anchor tag is missing or if the extracted value
// does not match the configured ID pattern.
func (u RawUnion) ID() (string, error) {
	val, exists := u.selection.Find("a.anchor").Attr("href")
	if !exists {
		return "", errors.New("attribute href not found")
	}
	if !u.idRegex.MatchString(val) {
		return "", fmt.Errorf("id %q does not match pattern %s", val, u.idRegex)
	}
	return val, nil
}

// Name returns the text content of the <h4> element.
//
// It returns an error if the content does not match the configured Name
// pattern.
func (u RawUnion) Name() (string, error) {
	val := u.selection.Text()
	if !u.nameRegex.MatchString(val) {
		return "", fmt.Errorf("name %q does not match pattern %s", val, u.nameRegex)
	}
	return val, nil
}

// Description extracts the documentation text for the union.
//
// It collects all <p> elements following the header up to the start of the
// variants list (<ul>) or the next section, joining them with spaces.
func (u RawUnion) Description() (string, error) {
	seq := u.selection.NextUntil("h1, h2, h3, h4, ul, table").Filter("p").All()
	var parts []string
	for _, p := range seq {
		parts = append(parts, p.Text())
	}
	return strings.Join(parts, " "), nil
}

// Variants returns an iterator over the possible objects (variants) of this
// union.
//
// It searches for the first <ul> list in the section and iterates over its <li>
// items.
func (u RawUnion) Variants() iter.Seq[Variant] {
	return func(yield func(Variant) bool) {
		seq := u.selection.NextUntil("h1, h2, h3, h4").Find("ul li").All()
		for _, li := range seq {
			if !yield(NewDefaultRawVariant(li)) {
				break
			}
		}
	}
}

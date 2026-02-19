// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"iter"
	"strings"

	"github.com/andreychh/tgen/parsing/dom"
)

// RawObject represents a lazily parsed Telegram API object definition.
//
// It expects the selection to point to a definition header, immediately
// followed by description paragraphs and a <table> of fields.
type RawObject struct {
	selection dom.Selection
}

// NewRawObject creates a RawObject starting from the given semantic DOM node.
func NewRawObject(h4 dom.Selection) RawObject {
	return RawObject{selection: h4}
}

// ID returns the unique reference identifier of the object (e.g., "#message").
func (t RawObject) ID() (string, error) {
	val, exists := t.selection.Find("a.anchor").Attr("href")
	if !exists {
		return "", errors.New("attribute href not found")
	}
	if !idRegex.MatchString(val) {
		return "", fmt.Errorf("id %q does not match pattern %s", val, idRegex)
	}
	return val, nil
}

// Name returns the name of the object.
func (t RawObject) Name() (string, error) {
	val := t.selection.Text()
	if !nameRegex.MatchString(val) {
		return "", fmt.Errorf("name %q does not match pattern %s", val, nameRegex)
	}
	return val, nil
}

// Description returns the human-readable context and purpose of the object.
func (t RawObject) Description() (string, error) {
	seq := t.selection.NextUntil("h1, h2, h3, h4, table").Filter("p").All()
	var parts []string
	for _, p := range seq {
		parts = append(parts, p.Text())
	}
	return strings.Join(parts, " "), nil
}

// Fields yields the properties defined for this object.
func (t RawObject) Fields() iter.Seq[Field] {
	return func(yield func(field Field) bool) {
		seq := t.selection.NextUntil("h1, h2, h3, h4").Find("table tbody tr").All()
		for _, tr := range seq {
			if !yield(NewRawObjectField(tr)) {
				break
			}
		}
	}
}

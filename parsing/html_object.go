// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"iter"
	"strings"

	"github.com/andreychh/tgen/parsing/dom"
)

// HTMLObject is an HTML-based implementation of the Object interface. It
// interprets a definition header (<h4>) in the Telegram Bot API documentation,
// followed by description paragraphs and an optional field table.
type HTMLObject struct {
	selection dom.Selection
}

// NewHTMLObject creates an HTMLObject from a definition header (<h4>) element.
func NewHTMLObject(h4 dom.Selection) HTMLObject {
	return HTMLObject{selection: h4}
}

// Ref returns the reference identifier of the object (e.g., "message").
func (o HTMLObject) Ref() (DefinitionRef, error) {
	val, exists := o.selection.Find("a.anchor").Attr("href")
	if !exists {
		return DefinitionRef{}, errors.New("anchor element is missing href attribute")
	}
	return NewDefinitionRef(strings.TrimPrefix(val, "#")), nil
}

// Name returns the name of the object (e.g., "Message").
func (o HTMLObject) Name() (ObjectName, error) {
	return NewObjectName(o.selection.Text()), nil
}

// Description returns the documentation text describing the object.
func (o HTMLObject) Description() (string, error) {
	seq := o.selection.NextUntil("h1, h2, h3, h4, table").Filter("p").All()
	var parts []string //nolint:prealloc // length is unknown before iteration
	for _, p := range seq {
		parts = append(parts, p.Text())
	}
	return strings.Join(parts, " "), nil
}

// Fields yields the properties defined for this object.
func (o HTMLObject) Fields() iter.Seq[Field] {
	return func(yield func(Field) bool) {
		seq := o.selection.NextUntil("h1, h2, h3, h4").Find("table tbody tr").All()
		for _, tr := range seq {
			if !yield(NewHTMLObjectField(tr)) {
				break
			}
		}
	}
}

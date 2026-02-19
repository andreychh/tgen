// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/andreychh/tgen/parsing/dom"
)

// DefinitionKind represents the category of a Telegram API definition.
type DefinitionKind string

const (
	// KindUnknown represents an unidentified or irrelevant documentation section.
	KindUnknown DefinitionKind = "unknown"

	// KindObject indicates that the definition is an Object.
	KindObject DefinitionKind = "object"

	// KindMethod indicates that the definition is a Method.
	KindMethod DefinitionKind = "method"

	// KindUnion indicates that the definition is a Union.
	KindUnion DefinitionKind = "union"
)

// Anchor represents a starting point of an API definition within the document.
//
// It wraps an <h4> header element and acts as a semantic marker to identify the
// kind of definition (e.g., [Method], [Object] or [Union]) before full parsing
// occurs.
type Anchor struct {
	selection dom.Selection
}

// NewAnchor creates a new Anchor instance from the provided <h4> DOM selection.
func NewAnchor(h4 dom.Selection) Anchor {
	return Anchor{selection: h4}
}

// Kind reports the category of the API definition pointed to by the anchor.
func (a Anchor) Kind() DefinitionKind {
	id, exists := a.selection.Find("a.anchor").Attr("href")
	if !exists || strings.Contains(id, "-") {
		return KindUnknown
	}
	first, _ := utf8.DecodeRuneInString(a.selection.Text())
	body := a.selection.NextUntil("h1, h2, h3, h4")
	hasTable := !body.Filter("table").IsEmpty()
	hasList := !body.Filter("ul").IsEmpty()
	switch {
	case unicode.IsLower(first):
		return KindMethod
	case unicode.IsUpper(first) && hasTable:
		return KindObject
	case unicode.IsUpper(first) && !hasTable && hasList:
		return KindUnion
	}
	return KindUnknown
}

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
	// KindUnknown represents a non-API or unidentified section.
	KindUnknown DefinitionKind = "unknown"

	// KindType represents a standard data object.
	KindType DefinitionKind = "type"

	// KindMethod represents an API endpoint.
	KindMethod DefinitionKind = "method"

	// KindUnion represents a sum-type interface.
	KindUnion DefinitionKind = "union"
)

// Anchor represents a starting point of an API definition within the document.
//
// It acts as a semantic marker used to identify the kind of entity (e.g. Method
// or Type) before full parsing occurs.
type Anchor struct {
	selection dom.Selection
}

// NewAnchor wraps a DOM selection to provide identification logic.
func NewAnchor(s dom.Selection) Anchor {
	return Anchor{selection: s}
}

// Kind identifies whether the anchor points to a Method, Type, or Union.
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
		return KindType
	case unicode.IsUpper(first) && !hasTable && hasList:
		return KindUnion
	}
	return KindUnknown
}

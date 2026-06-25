// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsed

import (
	"github.com/PuerkitoBio/goquery"
)

// Kind is the category of documentation section a Heading announces.
type Kind int

const (
	// KindUnknown marks a heading whose title matches no known category.
	KindUnknown Kind = iota
	KindObject
	KindMethod
	KindUnion
)

// Heading is the <h4> header of a documentation section. It tells which kind of
// section the header announces.
type Heading struct {
	h4 *goquery.Selection
}

// NewHeading constructs a Heading over an h4 selection.
func NewHeading(h4 *goquery.Selection) Heading {
	return Heading{h4: h4}
}

// Kind returns the category the heading announces, or KindUnknown when the
// title matches neither a type nor a method name.
func (h Heading) Kind() Kind {
	title := h.h4.Text()
	switch {
	case methodNamePattern.MatchString(title):
		return KindMethod
	case typeNamePattern.MatchString(title) && h.hasList():
		return KindUnion
	case typeNamePattern.MatchString(title):
		return KindObject
	}
	return KindUnknown
}

// hasList reports whether the section body holds a list, the source signal that
// a type enumerates variants and is therefore a union rather than a plain
// object.
func (h Heading) hasList() bool {
	return h.h4.NextUntil("h3, h4, hr").Filter("ul").Length() > 0
}

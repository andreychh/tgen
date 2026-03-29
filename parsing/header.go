// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/andreychh/tgen/parsing/gq"
)

type DefinitionKind string

const (
	DefinitionKindUnknown            DefinitionKind = "unknown"
	DefinitionKindObject             DefinitionKind = "object"
	DefinitionKindVariantObject      DefinitionKind = "variant_object"
	DefinitionKindMethod             DefinitionKind = "method"
	DefinitionKindUnion              DefinitionKind = "union"
	DefinitionKindDiscriminatedUnion DefinitionKind = "discriminated_union"
)

// GQHeader classifies an h4 section within a Telegram Bot API document. root is
// the document root used to resolve cross-document references when
// distinguishing discriminated unions from structural unions.
type GQHeader struct {
	root gq.Selection
	h4   gq.Selection
}

// NewGQHeader constructs a GQHeader from the document root and the h4 td of the
// definition.
func NewGQHeader(root, h4 gq.Selection) GQHeader {
	return GQHeader{root: root, h4: h4}
}

// Kind returns the kind of definition this h4 represents. For unions it
// performs a one-variant lookahead into the document to distinguish
// DefinitionKindDiscriminatedUnion from DefinitionKindUnion.
func (h GQHeader) Kind() DefinitionKind {
	id, exists := h.h4.Find("a.anchor").Attr("href")
	if !exists || strings.Contains(id, "-") {
		return DefinitionKindUnknown
	}
	first, _ := utf8.DecodeRuneInString(h.h4.Text())
	body := h.h4.Until("h3, h4, hr")
	hasList := !body.Filter("ul").IsEmpty()
	switch {
	case unicode.IsLower(first):
		return DefinitionKindMethod
	case unicode.IsUpper(first) && !hasList:
		return h.objectKind(body)
	case unicode.IsUpper(first) && hasList:
		return h.unionKind(body)
	}
	return DefinitionKindUnknown
}

func (h GQHeader) objectKind(body gq.Selection) DefinitionKind {
	hasDiscriminator := !body.
		Find("table tbody tr").
		FilterFunc(func(tr gq.Selection) bool {
			return NewGQFieldRow(tr).Kind() == FieldKindDiscriminator
		}).
		IsEmpty()
	if !hasDiscriminator {
		return DefinitionKindObject
	}
	name := h.h4.Text()
	isListed := !h.root.
		Find("div#dev_page_content h4").
		FilterFunc(func(cand gq.Selection) bool {
			return !cand.
				Until("h3, h4, hr").
				Find("ul li a").
				FilterFunc(func(a gq.Selection) bool {
					return a.Text() == name
				}).
				IsEmpty()
		}).
		IsEmpty()
	if isListed {
		return DefinitionKindVariantObject
	}
	return DefinitionKindObject
}

//nolint:varnamelen // <li> is the standard HTML list item element name
func (h GQHeader) unionKind(body gq.Selection) DefinitionKind {
	li := body.
		Find("ul li").
		At(0)
	if li.IsEmpty() {
		return DefinitionKindUnion
	}
	variant := h.root.
		Find("div#dev_page_content h4").
		FilterFunc(func(cand gq.Selection) bool {
			return cand.Text() == li.Find("a").Text()
		}).
		At(0)
	if variant.IsEmpty() {
		return DefinitionKindUnion
	}
	hasDiscriminator := !variant.
		Until("h3, h4, hr").
		Find("table tbody tr").
		FilterFunc(func(tr gq.Selection) bool {
			return NewGQFieldRow(tr).Kind() == FieldKindDiscriminator
		}).
		IsEmpty()
	if hasDiscriminator {
		return DefinitionKindDiscriminatedUnion
	}
	return DefinitionKindUnion
}

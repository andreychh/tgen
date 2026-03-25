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
	KindUnknown            DefinitionKind = "unknown"
	KindObject             DefinitionKind = "object"
	KindVariantObject      DefinitionKind = "variant_object"
	KindMethod             DefinitionKind = "method"
	KindUnion              DefinitionKind = "union"
	KindDiscriminatedUnion DefinitionKind = "discriminated_union"
)

// DefinitionHeader classifies an h4 section within a Telegram Bot API document.
// root is the document root used to resolve cross-document references when
// distinguishing discriminated unions from structural unions.
type DefinitionHeader struct {
	root gq.Selection
	h4   gq.Selection
}

// NewDefinitionHeader constructs a DefinitionHeader from the document root and
// the h4 selection of the definition.
func NewDefinitionHeader(root, h4 gq.Selection) DefinitionHeader {
	return DefinitionHeader{root: root, h4: h4}
}

// Kind returns the kind of definition this h4 represents. For unions it
// performs a one-variant lookahead into the document to distinguish
// KindDiscriminatedUnion from KindUnion.
func (h DefinitionHeader) Kind() DefinitionKind {
	id, exists := h.h4.Find("a.anchor").Attr("href")
	if !exists || strings.Contains(id, "-") {
		return KindUnknown
	}
	first, _ := utf8.DecodeRuneInString(h.h4.Text())
	body := h.h4.Until("h3, h4, hr")
	hasList := !body.Filter("ul").IsEmpty()
	switch {
	case unicode.IsLower(first):
		return KindMethod
	case unicode.IsUpper(first) && !hasList:
		for tr := range body.Find("table tbody tr").All() {
			if NewDefinitionRow(tr).Kind() == KindDiscriminatorField {
				return KindVariantObject
			}
		}
		return KindObject
	case unicode.IsUpper(first) && hasList:
		li := body.Find("ul li").At(0)
		if li.IsEmpty() {
			return KindUnion
		}
		variant := h.root.
			Find("div#dev_page_content h4").
			FilterFunc(func(cand gq.Selection) bool {
				return cand.Text() == li.Find("a").Text()
			}).
			At(0)
		if variant.IsEmpty() {
			return KindUnion
		}
		for tr := range variant.Until("h3, h4, hr").Find("table tbody tr").All() {
			if NewDefinitionRow(tr).Kind() == KindDiscriminatorField {
				return KindDiscriminatedUnion
			}
		}
		return KindUnion
	}
	return KindUnknown
}

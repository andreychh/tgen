// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/andreychh/tgen/pkg/gq"
)

type DefinitionKind string

const (
	DefinitionKindUnknown              DefinitionKind = "unknown"
	DefinitionKindObject               DefinitionKind = "object"
	DefinitionKindMethod               DefinitionKind = "method"
	DefinitionKindStructuredUnion      DefinitionKind = "structured_union"
	DefinitionKindDiscriminatedUnion   DefinitionKind = "discriminated_union"
	DefinitionKindFallbackUnion        DefinitionKind = "fallback_union"
	DefinitionKindDiscriminatedVariant DefinitionKind = "discriminated_variant"
)

// Header classifies an h4 section within a Telegram Bot API document. root is
// the document root used to resolve cross-document references when
// distinguishing union kinds.
type Header struct {
	root gq.Selection
	h4   gq.Selection
}

// NewHeader constructs a Header from the document root and the h4 td of the
// definition.
func NewHeader(root, h4 gq.Selection) Header {
	return Header{root: root, h4: h4}
}

// Kind returns the kind of definition this h4 represents. For unions it
// inspects all variants to distinguish DefinitionKindDiscriminatedUnion,
// DefinitionKindStructuredUnion, and DefinitionKindFallbackUnion.
func (h Header) Kind() DefinitionKind {
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

func (h Header) objectKind(body gq.Selection) DefinitionKind {
	hasDiscriminator := !body.
		Find("table tbody tr").
		FilterFunc(func(tr gq.Selection) bool {
			return NewFieldRow(tr).Kind() == FieldKindDiscriminator
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
		return DefinitionKindDiscriminatedVariant
	}
	return DefinitionKindObject
}

func (h Header) unionKind(body gq.Selection) DefinitionKind {
	all, discriminated := 0, 0
	for li := range body.Find("ul li").All() {
		variant := h.root.
			Find("div#dev_page_content h4").
			FilterFunc(func(cand gq.Selection) bool {
				return cand.Text() == li.Find("a").Text()
			}).
			At(0)
		if variant.IsEmpty() {
			continue
		}
		all++
		hasDiscriminator := !variant.
			Until("h3, h4, hr").
			Find("table tbody tr").
			FilterFunc(func(tr gq.Selection) bool {
				return NewFieldRow(tr).Kind() == FieldKindDiscriminator
			}).
			IsEmpty()
		if hasDiscriminator {
			discriminated++
		}
	}
	switch {
	case all == 0 || discriminated == 0:
		return DefinitionKindStructuredUnion
	case discriminated == all:
		return DefinitionKindDiscriminatedUnion
	default:
		return DefinitionKindFallbackUnion
	}
}

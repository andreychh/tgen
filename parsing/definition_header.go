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
	KindUnknown DefinitionKind = "unknown"
	KindObject  DefinitionKind = "object"
	KindMethod  DefinitionKind = "method"
	KindUnion   DefinitionKind = "union"
)

type DefinitionHeader struct {
	selection gq.Selection
}

func NewDefinitionHeader(h4 gq.Selection) DefinitionHeader {
	return DefinitionHeader{selection: h4}
}

func (h DefinitionHeader) Kind() DefinitionKind {
	id, exists := h.selection.Find("a.anchor").Attr("href")
	if !exists || strings.Contains(id, "-") {
		return KindUnknown
	}
	first, _ := utf8.DecodeRuneInString(h.selection.Text())
	hasList := !h.selection.
		Until("h3, h4, hr").
		Filter("ul").
		IsEmpty()
	switch {
	case unicode.IsLower(first):
		return KindMethod
	case unicode.IsUpper(first) && hasList:
		return KindUnion
	case unicode.IsUpper(first):
		return KindObject
	}
	return KindUnknown
}

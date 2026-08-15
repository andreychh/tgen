// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types_test

import (
	"strings"

	"github.com/andreychh/tgen/model/prose"
)

// plain builds the unemphasized run a method description writes its sentences
// in.
func plain(content string) prose.Text {
	return prose.NewText(content, prose.StylePlain)
}

// italic builds the emphasized run a return clause names a built-in type with.
func italic(content string) prose.Text {
	return prose.NewText(content, prose.StyleItalic)
}

// anchor builds the in-page link a return clause names a documented type with,
// addressing the section its name lowercases to.
func anchor(name string) prose.Link {
	return prose.NewLink(name, prose.StylePlain, "#"+strings.ToLower(name))
}

// passage builds a method description of one paragraph, the shape every return
// clause is read out of.
func passage(inlines ...prose.Inline) prose.Passage {
	return prose.NewPassage(prose.NewParagraph(inlines...))
}

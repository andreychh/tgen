// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package grammar_test

import (
	"strings"

	"github.com/andreychh/tgen/model/prose"
)

// plain builds the unemphasized run the documentation writes its sentences in.
func plain(content string) prose.Text {
	return prose.NewText(content, prose.StylePlain)
}

// italic builds the emphasized run the documentation sets a value apart with.
func italic(content string) prose.Text {
	return prose.NewText(content, prose.StyleItalic)
}

// anchor builds the in-page link the documentation names one of its own
// sections with, addressing the section the name lowercases to.
func anchor(name string) prose.Link {
	return prose.NewLink(name, prose.StylePlain, "#"+strings.ToLower(name))
}

// word is a rule matching a run that opens with one given plain text, standing
// in for the rules a pass writes.
type word struct {
	content string
	value   string
}

// newWord constructs a word matching content and decoding it into value.
func newWord(content, value string) word {
	return word{content: content, value: value}
}

// Match reports whether the run opens with the word's content and returns the
// value it stands for.
func (w word) Match(inlines []prose.Inline) (string, bool) {
	if len(inlines) == 0 {
		return "", false
	}
	text, ok := inlines[0].(prose.Text)
	if !ok || text.Content() != w.content {
		return "", false
	}
	return w.value, true
}

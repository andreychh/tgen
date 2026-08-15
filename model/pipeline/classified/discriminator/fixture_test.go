// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package discriminator_test

import (
	"strings"

	"github.com/andreychh/tgen/model/prose"
)

// plain builds the unemphasized run a field description writes its sentence in.
func plain(content string) prose.Text {
	return prose.NewText(content, prose.StylePlain)
}

// italic builds the emphasized run a must-be clause carries its value in.
func italic(content string) prose.Text {
	return prose.NewText(content, prose.StyleItalic)
}

// anchor builds the in-page link a field description names a documented type
// with, addressing the section its name lowercases to.
func anchor(name string) prose.Link {
	return prose.NewLink(name, prose.StylePlain, "#"+strings.ToLower(name))
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types_test

import (
	"strings"

	"github.com/andreychh/tgen/model/prose"
)

// plain builds the unemphasized run of text a type cell writes its keywords and
// built-in type words in.
func plain(content string) prose.Text {
	return prose.NewText(content, prose.StylePlain)
}

// anchor builds the in-page link a type cell names a documented type with,
// addressing the section its name lowercases to.
func anchor(name string) prose.Link {
	return prose.NewLink(name, prose.StylePlain, "#"+strings.ToLower(name))
}

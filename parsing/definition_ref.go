// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andreychh/tgen/parsing/gq"
)

var definitionRefRegex = regexp.MustCompile(`^#[a-z0-9]+$`)

type DefinitionRef struct {
	selection gq.Selection
}

func NewDefinitionRef(a gq.Selection) DefinitionRef {
	return DefinitionRef{selection: a}
}

func (r DefinitionRef) Value() (string, error) {
	href, _ := r.selection.Attr("href")
	if !definitionRefRegex.MatchString(href) {
		return "", fmt.Errorf("invalid definition ref: %q", href)
	}
	return strings.TrimPrefix(href, "#"), nil
}

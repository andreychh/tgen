// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/andreychh/tgen/parsing/gq"
)

var definitionRefRegex = regexp.MustCompile(`^#[a-z0-9]+$`)

type GQDefinitionReference struct {
	a gq.Selection
}

func NewGQDefinitionReference(a gq.Selection) GQDefinitionReference {
	return GQDefinitionReference{a: a}
}

func (r GQDefinitionReference) AsString() (string, error) {
	if r.a.IsEmpty() {
		return "", errors.New("definition ref not found")
	}
	href, _ := r.a.Attr("href")
	if !definitionRefRegex.MatchString(href) {
		return "", fmt.Errorf("definition ref %q contains invalid characters", href)
	}
	return strings.TrimPrefix(href, "#"), nil
}

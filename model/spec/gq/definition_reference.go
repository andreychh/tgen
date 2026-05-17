// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/pkg/gq"
)

var definitionRefRegex = regexp.MustCompile(`^#[a-z0-9]+$`)

type DefinitionReference struct {
	a gq.Selection
}

func NewDefinitionReference(a gq.Selection) DefinitionReference {
	return DefinitionReference{a: a}
}

// Value returns the definition reference extracted from the anchor href.
func (r DefinitionReference) Value() (model.Reference, error) {
	if r.a.IsEmpty() {
		return "", errors.New("definition ref not found")
	}
	href, _ := r.a.Attr("href")
	if !definitionRefRegex.MatchString(href) {
		return "", fmt.Errorf("definition ref %q contains invalid characters", href)
	}
	return model.Reference(strings.TrimPrefix(href, "#")), nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/pkg/gq"
)

var definitionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// Name represents an identifier parsed from the specification.
type Name struct {
	h4 gq.Selection
}

func NewName(h4 gq.Selection) Name {
	return Name{h4: h4}
}

func (n Name) AsString() (string, error) {
	if n.h4.IsEmpty() {
		return "", errors.New("definition name not found")
	}
	name := n.h4.Text()
	if !definitionNameRegex.MatchString(name) {
		return "", fmt.Errorf("definition name %q contains invalid characters", name)
	}
	return name, nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var definitionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// GQName represents an identifier parsed from the specification.
type GQName struct {
	h4 gq.Selection
}

func NewGQName(h4 gq.Selection) GQName {
	return GQName{h4: h4}
}

func (n GQName) AsString() (string, error) {
	if n.h4.IsEmpty() {
		return "", errors.New("definition name not found")
	}
	name := n.h4.Text()
	if !definitionNameRegex.MatchString(name) {
		return "", fmt.Errorf("definition name %q contains invalid characters", name)
	}
	return name, nil
}

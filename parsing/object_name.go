// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var objectNameRegex = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]+$`)

// GQObjectName represents a PascalCase identifier parsed from the specification.
type GQObjectName struct {
	selection gq.Selection
}

func NewGQObjectName(s gq.Selection) GQObjectName {
	return GQObjectName{selection: s}
}

func (n GQObjectName) Value() (string, error) {
	val := n.selection.Text()
	if !objectNameRegex.MatchString(val) {
		return "", fmt.Errorf("object name %q contains invalid characters", val)
	}
	return val, nil
}

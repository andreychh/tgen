// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var objectNameRegex = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]+$`)

type ObjectName struct {
	selection gq.Selection
}

func NewObjectName(s gq.Selection) ObjectName {
	return ObjectName{selection: s}
}

func (n ObjectName) Value() (string, error) {
	val := n.selection.Text()
	if !objectNameRegex.MatchString(val) {
		return "", fmt.Errorf("object name %q contains invalid characters", val)
	}
	return val, nil
}

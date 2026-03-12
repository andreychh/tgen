// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var methodNameRegex = regexp.MustCompile(`^[a-z][a-zA-Z0-9]+$`)

type MethodName struct {
	selection gq.Selection
}

func NewMethodName(s gq.Selection) MethodName {
	return MethodName{selection: s}
}

func (n MethodName) Value() (string, error) {
	val := n.selection.Text()
	if !methodNameRegex.MatchString(val) {
		return "", fmt.Errorf("method name %q contains invalid characters", val)
	}
	return val, nil
}

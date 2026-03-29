// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var keyRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

type GQKey struct {
	td gq.Selection
}

func NewGQKey(td gq.Selection) GQKey {
	return GQKey{td: td}
}

func (k GQKey) AsString() (string, error) {
	if k.td.IsEmpty() {
		return "", errors.New("field key not found")
	}
	key := k.td.Text()
	if !keyRegex.MatchString(key) {
		return "", fmt.Errorf("field key %q contains invalid characters", key)
	}
	return key, nil
}

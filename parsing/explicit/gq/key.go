// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/pkg/gq"
)

var keyRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

type Key struct {
	td gq.Selection
}

func NewKey(td gq.Selection) Key {
	return Key{td: td}
}

func (k Key) AsString() (string, error) {
	if k.td.IsEmpty() {
		return "", errors.New("field key not found")
	}
	key := k.td.Text()
	if !keyRegex.MatchString(key) {
		return "", fmt.Errorf("field key %q contains invalid characters", key)
	}
	return key, nil
}

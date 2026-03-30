// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/model"
)

type Tag struct {
	key         model.Key
	optionality model.Optionality
}

func NewTag(k model.Key, o model.Optionality) Tag {
	return Tag{key: k, optionality: o}
}

func (t Tag) AsString() (string, error) {
	key, err := t.key.AsString()
	if err != nil {
		return "", fmt.Errorf("getting field key: %w", err)
	}
	optional, err := t.optionality.AsBool()
	if err != nil {
		return "", fmt.Errorf("getting field optionality: %w", err)
	}
	if optional {
		return fmt.Sprintf("`json:\"%s,omitempty\"`", key), nil
	}
	return fmt.Sprintf("`json:%q`", key), nil
}

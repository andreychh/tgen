// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/parsing"
)

type FieldTag struct {
	key         parsing.FieldKey
	optionality parsing.FieldOptionality
}

func NewFieldTag(k parsing.FieldKey, o parsing.FieldOptionality) FieldTag {
	return FieldTag{key: k, optionality: o}
}

func (t FieldTag) Value() (string, error) {
	key, err := t.key.Value()
	if err != nil {
		return "", fmt.Errorf("getting field key: %w", err)
	}
	optional, err := t.optionality.Value()
	if err != nil {
		return "", fmt.Errorf("getting field optionality: %w", err)
	}
	if optional {
		return fmt.Sprintf("`json:\"%s,omitempty\"`", key), nil
	}
	return fmt.Sprintf("`json:\"%s\"`", key), nil
}

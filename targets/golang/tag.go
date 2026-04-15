// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/model"
)

type Tag struct {
	key model.Key
	opt model.Optionality
}

func NewTag(k model.Key, o model.Optionality) Tag {
	return Tag{key: k, opt: o}
}

func (t Tag) AsString() (string, error) {
	key, err := t.key.AsString()
	if err != nil {
		return "", fmt.Errorf("getting field key: %w", err)
	}
	opt, err := t.opt.AsBool()
	if err != nil {
		return "", fmt.Errorf("getting field optionality: %w", err)
	}
	if opt {
		return fmt.Sprintf("`json:\"%s,omitempty\"`", key), nil
	}
	return fmt.Sprintf("`json:%q`", key), nil
}

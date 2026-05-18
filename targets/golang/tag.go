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

func (t Tag) Value() string {
	if t.opt {
		return fmt.Sprintf("`json:\"%s,omitempty\"`", t.key)
	}
	return fmt.Sprintf("`json:%q`", t.key)
}

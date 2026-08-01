// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golangv2

import (
	"fmt"

	"github.com/andreychh/tgen/model"
)

// Tag represents the struct tag a field carries, telling encoding/json the key
// it reads and writes.
type Tag struct {
	key model.Key
	opt model.Optionality
}

// NewTag creates a Tag from a field key and its optionality.
func NewTag(key model.Key, opt model.Optionality) Tag {
	return Tag{key: key, opt: opt}
}

// NewRequiredTag creates a Tag for a field the encoded object always carries.
func NewRequiredTag(key model.Key) Tag {
	return NewTag(key, false)
}

// NewOptionalTag creates a Tag for a field the encoded object carries only
// when it holds something.
func NewOptionalTag(key model.Key) Tag {
	return NewTag(key, true)
}

// Value returns the tag, omitting an optional field from the encoded object
// when it carries nothing.
func (t Tag) Value() string {
	if t.opt {
		return fmt.Sprintf("`json:\"%s,omitempty\"`", t.key)
	}
	return fmt.Sprintf("`json:%q`", t.key)
}

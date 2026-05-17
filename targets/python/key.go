// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import "github.com/andreychh/tgen/model"

// Key represents a field key adapted for Python code generation.
type Key struct {
	inner model.Key
}

// NewKey constructs a Key from a parsed field key.
func NewKey(k model.Key) Key {
	return Key{inner: k}
}

func (k Key) Value() (string, error) {
	return string(k.inner), nil
}

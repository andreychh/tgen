// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package literals

type Key struct {
	value string
}

func NewKey(value string) Key {
	return Key{value: value}
}

func (k Key) AsString() (string, error) {
	return k.value, nil
}

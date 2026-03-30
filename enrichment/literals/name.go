// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package literals

// Name represents a parsing.Name wrapping a known string value.
type Name struct {
	value string
}

// NewName constructs a Name from value.
func NewName(value string) Name {
	return Name{value: value}
}

func (n Name) AsString() (string, error) {
	return n.value, nil
}

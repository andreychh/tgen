// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package literals

type DiscriminatorValue struct {
	value string
}

func NewDiscriminatorValue(value string) DiscriminatorValue {
	return DiscriminatorValue{value: value}
}

func (d DiscriminatorValue) AsString() (string, error) {
	return d.value, nil
}

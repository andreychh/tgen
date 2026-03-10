// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

type StaticName struct {
	origin string
}

func NewStaticName(s string) StaticName {
	return StaticName{origin: s}
}

func (n StaticName) Value() (string, error) {
	return n.origin, nil
}

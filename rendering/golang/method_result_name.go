// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

type MethodResultName struct {
	origin RawValue
}

func NewMethodResultName(n RawValue) MethodResultName {
	return MethodResultName{origin: n}
}

func (n MethodResultName) Value() (string, error) {
	val, err := n.origin.Value()
	if err != nil {
		return "", err
	}
	return val + "Result", nil
}

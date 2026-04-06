// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

type Array struct {
	element Expression
}

func NewArray(element Expression) Array {
	return Array{element: element}
}

func (a Array) Element() Expression {
	return a.element
}

func (a Array) Equals(other Expression) bool {
	if other, ok := other.(Array); ok {
		return a.element.Equals(other.element)
	}
	return false
}

func (a Array) isNode() {}

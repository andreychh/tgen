// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

// Array represents a homogeneous sequence of an element type.
type Array struct {
	element Expression
}

// NewArray constructs an array type from its element type.
func NewArray(element Expression) Array {
	return Array{element: element}
}

// Element returns the element type of the array.
func (a Array) Element() Expression {
	return a.element
}

// Equals reports whether other is an Array with an equal element type.
func (a Array) Equals(other Expression) bool {
	if other, ok := other.(Array); ok {
		return a.element.Equals(other.element)
	}
	return false
}

func (Array) isNode() {}

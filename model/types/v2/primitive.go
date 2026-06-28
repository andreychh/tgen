// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

// Primitive represents a built-in type the spec uses without an anchor.
type Primitive struct {
	kind PrimitiveKind
}

// NewPrimitive constructs a primitive type of the given kind.
func NewPrimitive(kind PrimitiveKind) Primitive {
	return Primitive{kind: kind}
}

// Kind returns the kind of the primitive.
func (p Primitive) Kind() PrimitiveKind {
	return p.kind
}

// Equals reports whether other is a Primitive of the same kind.
func (p Primitive) Equals(other Expression) bool {
	if other, ok := other.(Primitive); ok {
		return p.kind == other.kind
	}
	return false
}

func (Primitive) isNode() {}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package types defines TypeExpression, the algebraic type tree shared by the
// model and rendering layers. A TypeExpression is one of: a named type, an
// array type, or an inline union type.
package types

type TypeSource interface {
	AsString() (string, error)
}

// TypeExpression represents a node in the type expression tree parsed from the
// Telegram Bot API field table. A node is one of three kinds: a named type, an
// array type, or a union type.
//
// Exactly one of the three methods returns a non-zero value for any given node.
// The rendering layer traverses the tree by calling each method and acting on
// the one that returns true.
type TypeExpression interface {
	// Named returns the type name and true if this node is a named type (e.g.,
	// "Integer", "Message").
	Named() (string, bool)

	// Array returns the element type and true if this node is an array type (e.g.,
	// "Array of Integer").
	Array() (TypeExpression, bool)

	// Union returns the variant types and true if this node is an inline compound
	// type within a field definition (e.g., "Integer or String"). This is distinct
	// from a named polymorphic type (e.g., MaybeInaccessibleMessage) which is
	// defined as a top-level Union in the Specification.
	Union() ([]TypeExpression, bool)

	// Equals reports whether this node is structurally identical to other.
	Equals(other TypeExpression) bool
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import "github.com/andreychh/tgen/model"

// Named represents a reference to a documented type — an object or a union —
// addressed by the anchor of its definition.
type Named struct {
	ref model.Reference
}

// NewNamed constructs a named type from the reference of its definition.
func NewNamed(ref model.Reference) Named {
	return Named{ref: ref}
}

// Ref returns the reference of the named type's definition.
func (n Named) Ref() model.Reference {
	return n.ref
}

// Equals reports whether other is a Named with the same reference.
func (n Named) Equals(other Expression) bool {
	if other, ok := other.(Named); ok {
		return n.ref == other.ref
	}
	return false
}

func (Named) isNode() {}

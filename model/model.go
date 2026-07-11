// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package model defines the foundational value types shared across all layers
// of the tgen pipeline: Name, Type, Description, Key, and their companions.
package model

import "time"

type Name string

type Reference string

type Optionality bool

type Key string

// Position identifies a field's place among the other fields owned by the
// same reference: zero-based, unique and gapless within that owner's fields.
type Position int

type DiscriminatorValue string

type ReleaseVersion string

type ReleaseDate time.Time

// FieldKey identifies a field or parameter within the object or method that
// owns it, pairing the owner's reference with the field's own key.
type FieldKey struct {
	Owner Reference
	Key   Key
}

// VariantKey identifies a variant within the union that owns it, pairing the
// owner's reference with the variant's own reference.
type VariantKey struct {
	Owner Reference
	Ref   Reference
}

type Description interface {
	Value() (string, error)
	Links() ([]string, error)
}

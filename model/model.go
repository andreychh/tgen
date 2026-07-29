// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package model defines the foundational value types shared across all layers
// of the tgen pipeline: Name, Type, Description, Key, and their companions.
package model

import "time"

type Name string

type Reference string

// DefinitionKind classifies what a reference names, telling a target which
// shape to render it as. It is not the kind of a documentation heading: a
// heading may announce nothing recognizable, which no definition is, and no
// heading announces an alias, which tgen introduces on its own.
type DefinitionKind string

const (
	DefinitionKindObject DefinitionKind = "object"
	DefinitionKindMethod DefinitionKind = "method"
	DefinitionKindUnion  DefinitionKind = "union"
	DefinitionKindAlias  DefinitionKind = "alias"
)

type Optionality bool

type Key string

// Position is where a record's source put it: a field among the rows of its
// owner's table, a definition among the headings of the page. It is assigned
// once, when the record is read, and never recomputed, so a table that drops
// records keeps the gaps their positions leave behind. Which positions a
// given table admits is stated by that table.
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

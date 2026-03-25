// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"
)

//nolint:iface // DefinitionDescription and ObjectName are semantically distinct despite identical signatures
type DefinitionDescription interface {
	Value() (string, error)
}

//nolint:iface // ObjectName and DefinitionDescription are semantically distinct despite identical signatures
type ObjectName interface {
	Value() (string, error)
}

// Variant represents a named type in a union definition.
type Variant interface {
	Name() ObjectName
}

type Optionality interface {
	Value() (bool, error)
}

// Method represents a Telegram Bot API method definition.
type Method interface {
	Ref() DefinitionRef
	Name() MethodName
	Description() GQDefinitionDescription
	Returns() TypeTree
	Fields() iter.Seq[MethodField]
}

// FieldDescription represents the description column of a field table row.
type FieldDescription interface {
	Value() (string, error)
	// Links returns the hrefs of all anchor tags found in the description HTML.
	Links() []string
}

type Field interface {
	Key() FieldKey
	Type() TypeTree
	IsOptional() Optionality
	Description() FieldDescription
}

type TypeTree interface {
	Root() (TypeExpression, error)
}

// VariantObject represents an object that is a variant of a discriminated union.
// It is returned exclusively by DiscriminatedUnion.Variants() and never appears
// in Specification.Objects().
type VariantObject interface {
	Name() ObjectName
	Fields() VariantFields
	Ref() DefinitionRef
	Description() DefinitionDescription
}

// VariantFields groups the fields of a variant object: free fields that appear
// in the generated struct, and the discriminator field injected by MarshalJSON.
type VariantFields interface {
	Free() iter.Seq[Field]
	Discriminator() Discriminator
}

// Discriminator represents the key and fixed JSON value of a discriminator field
// for a specific variant (e.g. key="type", value="emoji").
type Discriminator interface {
	Key() FieldKey
	Value() (string, error)
}

type StructuralUnion interface {
	Ref() DefinitionRef
	Name() ObjectName
	Description() DefinitionDescription
	Variants() iter.Seq[VariantObject]
}

type DiscriminatedUnion interface {
	Ref() DefinitionRef
	Name() ObjectName
	Description() DefinitionDescription
	Variants() iter.Seq[VariantObject]
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package explicit defines the domain interfaces for data that comes directly
// from the Telegram Bot API specification: objects, methods, unions, fields,
// and their structural relationships. Concrete implementations live in the gq
// and overlays subpackages.
package explicit

import (
	"iter"

	"github.com/andreychh/tgen/model"
)

// Specification represents the full Telegram Bot API specification.
type Specification interface {
	Objects() iter.Seq[Object]
	Methods() iter.Seq[Method]
	Unions() Unions
	Release() Release
}

type Object interface {
	Reference() model.Reference
	Name() model.Name
	Description() model.Description
	Fields() iter.Seq[Field]
}

type Method interface {
	Reference() model.Reference
	Name() model.Name
	Description() model.Description
	ReturnType() model.Type
	Fields() iter.Seq[Field]
}

type Unions interface {
	Discriminated() iter.Seq[DiscriminatedUnion]
	Structured() iter.Seq[StructuredUnion]
}

type DiscriminatedUnion interface {
	Reference() model.Reference
	Name() model.Name
	Description() model.Description
	DiscriminatorKey() model.Key
	Variants() iter.Seq[DiscriminatedVariant]
}

type DiscriminatedVariant interface {
	Reference() model.Reference
	Name() model.Name
	Description() model.Description
	Fields() Fields
}

type Fields interface {
	Free() iter.Seq[Field]
	Discriminator() Discriminator
}

type Discriminator interface {
	Key() model.Key
	Value() model.DiscriminatorValue
}

type StructuredUnion interface {
	Reference() model.Reference
	Name() model.Name
	Description() model.Description
	Variants() iter.Seq[Object]
}

type Field interface {
	Key() model.Key
	Type() model.Type
	Optionality() model.Optionality
	Description() model.Description
}

type Release interface {
	Reference() model.Reference
	Version() model.ReleaseVersion
	Date() model.ReleaseDate
}

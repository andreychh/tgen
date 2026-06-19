// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package spec defines the domain interfaces for data that comes directly
// from the Telegram Bot API specification: objects, methods, unions, fields,
// and their structural relationships. Concrete implementations live in the gq
// and overlays subpackages.
package spec

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/types"
)

// Specification represents the full Telegram Bot API specification.
type Specification interface {
	Objects() iter.Seq[Object]
	Methods() iter.Seq[Method]
	DiscriminatedObjects() iter.Seq[DiscriminatedObject]
	DiscriminatedUnions() iter.Seq[DiscriminatedUnion]
	Release() Release
}

type Object interface {
	Reference() (model.Reference, error)
	Name() (model.Name, error)
	Description() model.Description
	Fields() iter.Seq[Field]
}

type Method interface {
	Reference() (model.Reference, error)
	Name() (model.Name, error)
	Description() model.Description
	Result() (Result, error)
	Fields() iter.Seq[Field]
}

type DiscriminatedUnion interface {
	Reference() (model.Reference, error)
	Name() (model.Name, error)
	Description() model.Description
	DiscriminatorKey() (model.Key, error)
	Variants() iter.Seq[DiscriminatedObject]
}

type DiscriminatedObject interface {
	Reference() (model.Reference, error)
	Name() (model.Name, error)
	Description() model.Description
	Fields() Fields
}

type Fields interface {
	Free() iter.Seq[Field]
	Discriminator() Discriminator
}

type Discriminator interface {
	Key() (model.Key, error)
	Value() (model.DiscriminatorValue, error)
}

type Field interface {
	Key() (model.Key, error)
	Type() (types.Expression, error)
	Optionality() (model.Optionality, error)
	Description() model.Description
}

type Release interface {
	Reference() (model.Reference, error)
	Version() (model.ReleaseVersion, error)
	Date() (model.ReleaseDate, error)
}

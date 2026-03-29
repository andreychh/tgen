// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"iter"
	"time"

	"github.com/andreychh/tgen/parsing/types"
)

// Specification represents the full Telegram Bot API specification.
type Specification interface {
	Objects() iter.Seq[Object]
	Methods() iter.Seq[Method]
	Unions() Unions
	Release() Release
}

type Object interface {
	Reference() Reference
	Name() Name
	Description() Description
	Fields() iter.Seq[Field]
}

type Method interface {
	Reference() Reference
	Name() Name
	Description() Description
	ReturnType() Type
	Fields() iter.Seq[Field]
}

type Unions interface {
	Discriminated() iter.Seq[DiscriminatedUnion]
	Structured() iter.Seq[StructuredUnion]
}

type DiscriminatedUnion interface {
	Reference() Reference
	Name() Name
	Description() Description
	DiscriminatorKey() Key
	Variants() iter.Seq[DiscriminatedVariant]
}

type DiscriminatedVariant interface {
	Reference() Reference
	Name() Name
	Description() Description
	Fields() Fields
}

type Fields interface {
	Free() iter.Seq[Field]
	Discriminator() Discriminator
}

type Discriminator interface {
	Key() Key
	Value() DiscriminatorValue
}

type StructuredUnion interface {
	Reference() Reference
	Name() Name
	Description() Description
	Variants() iter.Seq[StructuredVariant]
}

type StructuredVariant interface {
	Reference() Reference
	Name() Name
	Description() Description
	Fields() iter.Seq[Field]
}

type Field interface {
	Key() Key
	Type() Type
	Optionality() Optionality
	Description() Description
}

type Release interface {
	Reference() Reference
	Version() Version
	Date() Date
}

// -----

type Name interface {
	AsString() (string, error)
}

type Reference interface {
	AsString() (string, error)
}

type Type interface {
	AsExpression() (types.TypeExpression, error)
}

type Optionality interface {
	AsBool() (bool, error)
}

type Description interface {
	AsString() (string, error)
	Links() ([]string, error)
}

type Key interface {
	AsString() (string, error)
}

type DiscriminatorValue interface {
	AsString() (string, error)
}

type Version interface {
	AsString() (string, error)
}

type Date interface {
	AsTime() (time.Time, error)
}

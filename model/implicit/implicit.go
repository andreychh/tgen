// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package implicit defines types invented by tgen that have no counterpart in
// the Telegram Bot API specification, such as discriminated unions constructed
// from patterns tgen recognises across the explicit layer.
package implicit

import (
	"iter"

	"github.com/andreychh/tgen/model"
)

type Specification interface {
	Unions() Unions
}

// Unions represents the set of implicitly defined union types.
type Unions interface {
	Discriminated() iter.Seq[DiscriminatedUnion]
	Structured() iter.Seq[StructuredUnion]
}

// DiscriminatedUnion represents an implicitly defined union whose variants are
// distinguished by a discriminator field value.
type DiscriminatedUnion interface {
	Name() model.Name
	Description() model.Description
	DiscriminatorKey() model.Key
	Variants() iter.Seq[DiscriminatedVariant]
}

// DiscriminatedVariant represents a named variant of a DiscriminatedUnion.
type DiscriminatedVariant interface {
	Name() model.Name
	DiscriminatorValue() model.DiscriminatorValue
}

// StructuredUnion represents an implicitly defined union whose variants are
// distinguished by structure.
type StructuredUnion interface {
	Name() model.Name
	Description() model.Description
	Variants() iter.Seq[StructuredVariant]
}

// StructuredVariant represents a named variant of a StructuredUnion with a
// distinct field type.
type StructuredVariant interface {
	Name() model.Name
	Type() model.Type
}

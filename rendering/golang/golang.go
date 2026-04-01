// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package golang provides the necessary templates and execution context to
// generate Go source code from the parsed Telegram Bot API specification.
package golang

import "iter"

const specificationURL = "https://core.telegram.org/bots/api"

type DiscriminatedUnion interface {
	Name() Name
	Doc() GoDoc
	DiscriminatorKey() Key
	Variants() iter.Seq[DiscriminatedVariant]
}

type DiscriminatedVariant interface {
	Name() Name
	DiscriminatorValue() DiscriminatorValue
}

type StructuredUnion interface {
	Name() Name
	Doc() GoDoc
	Variants() iter.Seq[StructuredVariant]
}

type StructuredVariant interface {
	Name() Name
	Type() Type
}

// Type represents a Go type expression that can be rendered as a string.
//
//nolint:iface // intentionally distinct from Stringable: Type is a Go type expression, not a doc string
type Type interface {
	AsString() (string, error)
}

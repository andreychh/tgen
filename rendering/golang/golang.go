// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package golang provides the necessary templates and execution context to
// generate Go source code from the parsed Telegram Bot API specification.
package golang

import "iter"

const specificationURL = "https://core.telegram.org/bots/api"

type Object interface {
	Name() Name
	Doc() GoDoc
	Fields() iter.Seq[Field]
}

type DiscriminatedUnion interface {
	Name() Name
	Doc() GoDoc
	DiscriminatorKey() Key
	Variants() iter.Seq[DiscriminatedVariant]
}

type DiscriminatedVariant interface {
	Name() Name
	Type() Name
	Doc() GoDoc
	Fields() DiscriminatedVariantFields
}

type StructuredUnion interface {
	Name() Name
	Doc() GoDoc
	Variants() iter.Seq[Object]
}

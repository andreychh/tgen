// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package parsing provides the core logic for extracting structured API
// definitions from the Telegram Bot API HTML documentation.
//
// It defines the domain models (e.g. Document, Union, Variant) and their
// HTML-based implementations, serving as the data source for the code
// generator.
package parsing

import (
	"iter"
	"regexp"
)

var (
	// idRegex matches valid anchor identifiers used in the documentation (e.g.,
	// "#message", "#getupdates").
	idRegex = regexp.MustCompile(`^#[a-z0-9]+$`)

	// nameRegex matches valid type names in PascalCase (e.g., "Message", "User").
	nameRegex = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]+$`)
)

// Document represents the Telegram API specification.
type Document interface {
	// Unions returns a sequence of all Union types found in the specification.
	Unions() iter.Seq[Union]
}

// Union represents a polymorphic type definition (Sum Type) where a value can
// be one of several distinct types.
//
// Example: In the Telegram Bot API, "MaybeInaccessibleMessage" is a union of
// "Message" and "InaccessibleMessage".
type Union interface {
	// ID returns the unique reference identifier of the union definition (e.g.
	// "#maybeinaccessiblemessage").
	ID() (string, error)

	// Name returns the type name of the union (e.g. "MaybeInaccessibleMessage").
	Name() (string, error)

	// Description returns the documentation text associated with the union.
	Description() (string, error)

	// Variants returns an iterator over the possible types that form this union.
	Variants() iter.Seq[Variant]
}

// Variant represents a single option within a Union type definition.
//
// Example: In the context of Telegram Bot API, the "MaybeInaccessibleMessage"
// union contains two variants: "Message" and "InaccessibleMessage".
type Variant interface {
	// ID returns the unique reference identifier of the variant (e.g. "#message").
	ID() (string, error)

	// Name returns the type name of the variant (e.g. "Message").
	Name() (string, error)
}

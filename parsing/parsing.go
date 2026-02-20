// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package parsing provides the core logic for extracting structured API
// definitions from the Telegram Bot API HTML documentation.
//
// It defines the domain models (e.g., Specification, Union, Variant) and their
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

	// nameRegex matches valid Object and Union names in PascalCase (e.g.,
	// "Message", "User").
	nameRegex = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]+$`)

	// jsonKeyRegex matches valid JSON property keys (e.g., "message_id").
	jsonKeyRegex = regexp.MustCompile(`^[a-z_]+$`)

	// typeRegex matches valid Field types, which can include basic types (e.g.,
	// "Integer", "String") or array types (e.g., "Array of Integer").
	typeRegex = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
)

// Specification represents the contract of the Telegram Bot API.
type Specification interface {
	// Unions returns a sequence of all union definitions found in the API.
	Unions() iter.Seq[Union]

	// Objects returns a sequence of all object definitions found in the API.
	Objects() iter.Seq[Object]
}

// Object represents a standard object definition in the Telegram Bot API. It
// defines the structure of the data payloads sent to and received from the
// Telegram.
//
// Example: In the Telegram Bot API, "Message" and "User" are objects.
type Object interface {
	// ID returns the unique reference identifier of the object definition (e.g.,
	// "#message").
	ID() (string, error)

	// Name returns the name of the object, formatted to match a specific code
	// convention (e.g., "Message").
	Name() (string, error)

	// Description returns the documentation text describing the object.
	Description() (string, error)

	// Fields returns an iterator over the properties that belong to this object.
	Fields() iter.Seq[Field]
}

// Field represents a single property within an Object definition.
type Field interface {
	// Name returns the name of the field, formatted to match a specific code
	// convention (e.g., "MessageID").
	Name() (string, error)

	// Description returns the documentation text explaining the field's purpose.
	Description() (string, error)

	// Type returns the data type of the field, formatted to match a specific code
	// convention (e.g., "int", "string", or "[]User").
	Type() (string, error)

	// JSONKey returns the exact string key used to serialize and deserialize this
	// field in JSON payloads over the network (e.g., "message_id").
	JSONKey() (string, error)

	// IsOptional reports whether the field is not strictly required in the JSON
	// payload.
	IsOptional() (bool, error)
}

// Union represents a polymorphic definition (Sum Type) where a value can be one
// of several distinct objects.
//
// Example: In the Telegram Bot API, "MaybeInaccessibleMessage" is a union of
// "Message" and "InaccessibleMessage".
type Union interface {
	// ID returns the unique reference identifier of the union definition (e.g.,
	// "#maybeinaccessiblemessage").
	ID() (string, error)

	// Name returns the name of the union, formatted to match a specific code
	// convention (e.g., "MaybeInaccessibleMessage").
	Name() (string, error)

	// Description returns the documentation text associated with the union.
	Description() (string, error)

	// Variants returns an iterator over the possible objects that form this union.
	Variants() iter.Seq[Variant]
}

// Variant represents a single option within a Union definition.
//
// Example: In the Telegram Bot API, the "MaybeInaccessibleMessage" union
// contains two variants: "Message" and "InaccessibleMessage".
type Variant interface {
	// ID returns the unique reference identifier of the variant (e.g., "#message").
	ID() (string, error)

	// Name returns the name of the variant, formatted to match a specific code
	// convention (e.g., "Message").
	Name() (string, error)
}

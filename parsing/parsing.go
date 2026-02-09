// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package parsing provides the core logic for extracting structured API
// definitions from the Telegram Bot API HTML documentation.
//
// It defines the domain models (Definition, Union, Variant) and their
// HTML-based implementations, serving as the data source for the code
// generator.
package parsing

import (
	"regexp"
)

var (
	// idRegex matches valid anchor identifiers used in the documentation (e.g.,
	// "#message", "#getupdates").
	idRegex = regexp.MustCompile(`^#[a-z0-9]+$`)

	// nameRegex matches valid type names in PascalCase (e.g., "Message", "User").
	nameRegex = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]+$`)
)

// Variant represents a single option within a Union type definition.
//
// Example: In the context of Telegram Bot API, the "MaybeInaccessibleMessage"
// union contains two variants: "Message" and "InaccessibleMessage".
type Variant interface {
	// ID returns the unique reference identifier of the variant.
	//
	// In HTML documentation, this typically corresponds to the anchor href (e.g.,
	// "#message").
	ID() (string, error)

	// Name returns the type name of the variant.
	//
	// Implementations MUST perform whitespace normalization: trimming leading and
	// trailing spaces to ensure the returned string is a valid type identifier.
	Name() (string, error)
}

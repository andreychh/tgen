// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

type Optionality interface {
	Value() (bool, error)
}

// FieldDescription represents the description column of a field table row.
type FieldDescription interface {
	RawValue
	// Links returns the hrefs of all anchor tags found in the description HTML.
	Links() []string
}

type Field interface {
	Key() FieldKey
	Type() TypeTree
	IsOptional() Optionality
	Description() FieldDescription
}

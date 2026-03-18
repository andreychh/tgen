// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import "iter"

type DefinitionDescription interface {
	Value() (string, error)
}

type ObjectName interface {
	Value() (string, error)
}

// Variant represents a named type in a union definition.
type Variant interface {
	Name() ObjectName
}

type Optionality interface {
	Value() (bool, error)
}

// Method represents a Telegram Bot API method definition.
type Method interface {
	Ref() DefinitionRef
	Name() MethodName
	Description() GQDefinitionDescription
	Returns() TypeTree
	Fields() iter.Seq[MethodField]
}

// FieldDescription represents the description column of a field table row.
type FieldDescription interface {
	Value() (string, error)
	// Links returns the hrefs of all anchor tags found in the description HTML.
	Links() []string
}

type Field interface {
	Key() FieldKey
	Type() TypeTree
	IsOptional() Optionality
	Description() FieldDescription
}

type TypeTree interface {
	Root() (TypeExpression, error)
}

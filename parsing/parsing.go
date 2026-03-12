// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

type Optionality interface {
	Value() (bool, error)
}

type Field interface {
	Key() FieldKey
	Type() TypeTree
	IsOptional() Optionality
	Description() RawValue
}

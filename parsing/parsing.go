// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

type Field interface {
	Key() FieldKey
	Type() TypeTree
	IsOptional() FieldOptionality
	Description() FieldDescription
}

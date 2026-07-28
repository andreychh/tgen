// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package primitive enumerates the built-in types the Telegram Bot API
// documentation writes without an anchor of their own.
package primitive

// Kind identifies one built-in type.
type Kind string

const (
	Integer Kind = "Integer"
	String  Kind = "String"
	Boolean Kind = "Boolean"
	Float   Kind = "Float"
	True    Kind = "True"
)

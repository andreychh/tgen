// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package golang provides the necessary templates and execution context to
// generate Go source code from the parsed Telegram Bot API specification.
package golang

type RawValue interface {
	Value() (string, error)
}

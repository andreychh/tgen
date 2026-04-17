// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

// Type represents a Go type expression used in generated code.
type Type interface {
	IsUnion() (bool, error)
	IsPrimitive() (bool, error)
	Depth() (int, error)
	Name() (string, error)
	AsString() (string, error)
	Zero() (string, error)
	Part() (string, error)
}

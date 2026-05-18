// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

// Shape categorizes a Type by how the Go code generator must handle it.
type Shape string

const (
	// ShapeUnionArray marks a slice of a discriminated union.
	ShapeUnionArray Shape = "union_array"
	// ShapeUnion marks a scalar discriminated union.
	ShapeUnion Shape = "union"
	// ShapePlain marks any other type rendered directly.
	ShapePlain Shape = "plain"
)

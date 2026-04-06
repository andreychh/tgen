// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

type Kind string

const (
	KindUnknown   Kind = ""
	KindObject    Kind = "object"
	KindUnion     Kind = "union"
	KindPrimitive Kind = "primitive"
)

type Catalog interface {
	Lookup(name string) (Kind, bool)
}

//sumtype:decl
type Expression interface {
	Equals(other Expression) bool
	isNode()
}

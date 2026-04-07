// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

type Named struct {
	name string
	kind Kind
}

func NewNamed(name string, k Kind) Named {
	return Named{name: name, kind: k}
}

func (n Named) Name() string {
	return n.name
}

func (n Named) Kind() Kind {
	return n.kind
}

func (n Named) Equals(other Expression) bool {
	if other, ok := other.(Named); ok {
		return n.name == other.name && n.kind == other.kind
	}
	return false
}

func (n Named) String() string {
	return n.name + "(" + string(n.kind) + ")"
}

func (n Named) isNode() {}

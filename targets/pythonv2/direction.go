// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import "github.com/andreychh/tgen/model"

// Direction is which way a declaration travels, read as what has to be written
// for it. The three predicates name the exact direction, and every declaration
// answers yes to one of them.
type Direction struct {
	inner model.Direction
}

// NewDirection creates a Direction from the way a declaration travels.
func NewDirection(direction model.Direction) Direction {
	return Direction{inner: direction}
}

// Outbound reports whether a request alone carries the declaration.
func (d Direction) Outbound() bool {
	return d.inner == model.DirectionOutbound
}

// Inbound reports whether a response alone carries the declaration.
func (d Direction) Inbound() bool {
	return d.inner == model.DirectionInbound
}

// Bidirectional reports whether a request and a response both carry the
// declaration.
func (d Direction) Bidirectional() bool {
	return d.inner == model.DirectionBidirectional
}

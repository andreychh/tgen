// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/model"

// Direction is which way a declaration travels, read as what has to be written
// for it. A declaration travelling one way needs the half of the codec that way
// calls for and no more, since nothing encodes a value no request ever carries
// and nothing decodes one no response ever carries.
//
// Two kinds of question are asked of it and their answers must not be confused.
// [Direction.Sent] and [Direction.Received] report what a declaration is capable
// of, and one travelling both ways answers yes to each. [Direction.Outbound],
// [Direction.Inbound] and [Direction.Bidirectional] name the exact direction,
// and every declaration answers yes to one of the three.
type Direction struct {
	inner model.Direction
}

// NewDirection creates a Direction from the way a declaration travels.
func NewDirection(direction model.Direction) Direction {
	return Direction{inner: direction}
}

// Sent reports whether a request ever carries the declaration, which is what
// obliges it to write itself into JSON.
func (d Direction) Sent() bool {
	return d.Outbound() || d.Bidirectional()
}

// Received reports whether a response ever carries the declaration, which is
// what obliges it to read itself out of JSON.
func (d Direction) Received() bool {
	return d.Inbound() || d.Bidirectional()
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

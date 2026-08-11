// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import "github.com/andreychh/tgen/model"

// Direction is which way a declaration travels, read as what has to be written
// for it. A declaration travelling one way needs the half of the codec that way
// calls for and no more, since nothing encodes a value no request ever carries
// and nothing decodes one no response ever carries.
type Direction struct {
	inner model.Direction
}

// NewDirection creates a Direction from the way a declaration travels.
func NewDirection(direction model.Direction) Direction {
	return Direction{inner: direction}
}

// Sent reports whether a request ever carries the declaration, which is what
// obliges it to write itself into JSON. One travelling both ways is sent no
// less than one travelling outbound alone.
func (d Direction) Sent() bool {
	return d.inner == model.DirectionOutbound || d.inner == model.DirectionBidirectional
}

// Received reports whether a response ever carries the declaration, which is
// what obliges it to read itself out of JSON. One travelling both ways is
// received no less than one travelling inbound alone.
func (d Direction) Received() bool {
	return d.inner == model.DirectionInbound || d.inner == model.DirectionBidirectional
}

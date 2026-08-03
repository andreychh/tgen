// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package directed

import (
	"errors"
	"fmt"

	"github.com/andreychh/tgen/model"
)

// Reach is how far the spread has got to one definition. It is the state of a
// traversal and not a property of what is traversed, which is why it holds a
// fourth case [model.Direction] has no room for: a definition may be reached by
// nothing yet, and while the spread runs most of them are.
type Reach string

const (
	// ReachNone marks a definition nothing has reached. It is the zero value on
	// purpose, so a reference absent from the spread reads as unreached without
	// being asked for separately.
	ReachNone Reach = ""
	// ReachOutbound marks a definition reached only as something a request
	// carries.
	ReachOutbound Reach = "outbound"
	// ReachInbound marks a definition reached only as something a response
	// carries.
	ReachInbound Reach = "inbound"
	// ReachBoth marks a definition reached as both.
	ReachBoth Reach = "both"
)

// Merged returns what is known of a definition reached both as r and as other.
// The answer is never less than either, which is what settles the spread: a
// reach rises twice at most, from [ReachNone] to one direction and from there
// to [ReachBoth].
func (r Reach) Merged(other Reach) Reach {
	switch {
	case r == other:
		return r
	case r == ReachNone:
		return other
	case other == ReachNone:
		return r
	default:
		return ReachBoth
	}
}

// Direction returns the direction the reach amounts to. It fails when nothing
// reached the definition, a definition no request and no response carries being
// a gap in the specification rather than a direction of its own, and when the
// reach is one no spread produces.
func (r Reach) Direction() (model.Direction, error) {
	switch r {
	case ReachOutbound:
		return model.DirectionOutbound, nil
	case ReachInbound:
		return model.DirectionInbound, nil
	case ReachBoth:
		return model.DirectionBidirectional, nil
	case ReachNone:
		return "", errors.New("reached by nothing")
	default:
		return "", fmt.Errorf("unknown reach %q", r)
	}
}

// Reaches is the state of the spread: how far every definition has been reached
// so far, keyed by reference. It is the one thing a spread mutates, and it
// mutates only upwards. Its zero value is not usable; construct one with
// [NewReaches]. A Reaches is not safe for concurrent use.
type Reaches struct {
	inner map[model.Reference]Reach
}

// NewReaches constructs an empty Reaches.
func NewReaches() Reaches {
	return Reaches{inner: make(map[model.Reference]Reach)}
}

// Lookup returns how far ref has been reached, [ReachNone] for a reference
// nothing has reached yet.
func (r Reaches) Lookup(ref model.Reference) Reach {
	return r.inner[ref]
}

// Merge records that ref is reached as reach besides however it was reached
// before, and reports whether that raised what is known of it. Merging
// [ReachNone] raises nothing and records nothing, which is what lets an edge
// leading out of an unreached definition be followed without a test of its own.
func (r Reaches) Merge(ref model.Reference, reach Reach) bool {
	known := r.inner[ref]
	merged := known.Merged(reach)
	if merged == known {
		return false
	}
	r.inner[ref] = merged
	return true
}

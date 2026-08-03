// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package directed

import "github.com/andreychh/tgen/model"

// Seed is one definition the traversal starts from, and how it is reached. A
// seed stands for what no [Rule] can say: a parameter goes out because it is
// part of a request, not because anything reached the method taking it, and a
// relation producing a reach instead of passing one along is no edge of the
// graph.
type Seed struct {
	Ref   model.Reference
	Reach Reach
}

// Spread is the traversal: it starts from the seeds and follows every edge the
// rules put in the graph until a round raises nothing.
type Spread struct {
	seeds []Seed
	rules []Rule
}

// NewSpread constructs a Spread starting from seeds and following the edges of
// the given rules.
func NewSpread(seeds []Seed, rules ...Rule) Spread {
	return Spread{seeds: seeds, rules: rules}
}

// Value returns how far every definition was reached. The rounds repeat until
// one raises nothing, so a definition holding its own kind — a rich block
// nesting rich blocks — settles instead of circling.
func (s Spread) Value() Reaches {
	out := NewReaches()
	for _, seed := range s.seeds {
		out.Merge(seed.Ref, seed.Reach)
	}
	for s.round(out) {
	}
	return out
}

// round follows every edge once, carrying what is known of the source into the
// target, and reports whether that raised anything. An edge leading out of a
// definition nothing reached carries [ReachNone] and raises nothing, so no edge
// needs a test of its own.
func (s Spread) round(out Reaches) bool {
	grew := false
	for _, rule := range s.rules {
		for from, to := range rule.Edges() {
			raised := out.Merge(to, out.Lookup(from))
			grew = grew || raised
		}
	}
	return grew
}

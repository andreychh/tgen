// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package directed

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/flattened"
	"github.com/andreychh/tgen/model/typeform"
)

// envelopeRef is the reference of the object naming why a request failed, which
// a response carries beside the result of whichever method was called. No method
// names it, so no edge of the graph leads to it: a client decodes it because
// every response has room for it, and that is knowledge of the transport rather
// than of any one method.
const envelopeRef = model.Reference("responseparameters")

// Seeding is where the traversal starts: what a transport touches on its own,
// without anything having had to reach it. A request carries the parameters of
// the method it invokes, a response carries what that method hands back, and
// either way a definition reached this way is reached because of what it is
// part of, not because of what reached it.
type Seeding struct {
	methods flattened.Methods
	fields  flattened.Fields
}

// NewSeeding constructs a Seeding over the tables of methods and fields.
func NewSeeding(methods flattened.Methods, fields flattened.Fields) Seeding {
	return Seeding{methods: methods, fields: fields}
}

// Value returns every definition the traversal starts from and how each of them
// is reached. A reference may be named more than once, one method taking it as
// a parameter while another hands it back, and it is the merging of the seeds
// that leaves such a definition travelling both ways.
func (s Seeding) Value() []Seed {
	return append(s.outbound(), s.inbound()...)
}

// outbound returns one seed for every parameter of a method naming a
// definition, since a request carries whatever the method it invokes was given.
// A parameter typed as a primitive names no definition and seeds nothing, and a
// field of an object is no parameter, however alike the two are held.
func (s Seeding) outbound() []Seed {
	out := make([]Seed, 0)
	for key, field := range s.fields.All() {
		if _, invoked := s.methods.Lookup(key.Owner); !invoked {
			continue
		}
		named, ok := field.Type.Atom().(typeform.Named)
		if !ok {
			continue
		}
		out = append(out, Seed{Ref: named.Ref(), Reach: ReachOutbound})
	}
	return out
}

// inbound returns one seed for every method return naming a definition, and one
// for the envelope every response carries besides. A method handing back a
// primitive — the bare True of a command — names no definition and seeds
// nothing.
func (s Seeding) inbound() []Seed {
	out := []Seed{{Ref: envelopeRef, Reach: ReachInbound}}
	for _, method := range s.methods.All() {
		named, ok := method.Type.Atom().(typeform.Named)
		if !ok {
			continue
		}
		out = append(out, Seed{Ref: named.Ref(), Reach: ReachInbound})
	}
	return out
}

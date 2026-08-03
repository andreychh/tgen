// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package directed

import (
	"iter"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/flattened"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/typeform"
)

// Rule is one relation of the specification read as edges of the graph a
// direction travels along. An edge says one thing and only that: the two
// travel together, so whatever reaches the source reaches the target
// unchanged. That is why a rule names no direction of its own — it knows which
// definitions go together, and what they carry is decided where the marks are
// kept. A relation obliged to produce a direction instead of passing one along
// is no rule at all, and enters the spread as a seed.
type Rule interface {
	// Edges returns every pair the relation puts in the graph, the source first
	// and the target second, in no particular order. A relation reaching a
	// primitive puts no pair in the graph, a primitive being no definition an
	// edge could lead to.
	Edges() iter.Seq2[model.Reference, model.Reference]
}

// FieldRule is the [Rule] of ownership: an object travels whole, so the type of
// a field goes wherever the definition owning that field goes. Dimensionality
// does not matter, since an array travels with what it holds. A field typed as
// a primitive puts no edge in the graph, a primitive being no definition that
// could travel anywhere.
type FieldRule struct {
	fields flattened.Fields
}

// NewFieldRule constructs a FieldRule over a table of fields.
func NewFieldRule(fields flattened.Fields) FieldRule {
	return FieldRule{fields: fields}
}

// Edges implements [Rule].
func (r FieldRule) Edges() iter.Seq2[model.Reference, model.Reference] {
	return func(yield func(model.Reference, model.Reference) bool) {
		for key, field := range r.fields.All() {
			named, ok := field.Type.Atom().(typeform.Named)
			if !ok {
				continue
			}
			if !yield(key.Owner, named.Ref()) {
				return
			}
		}
	}
}

// VariantRule is the [Rule] of alternation: a union travels as whichever
// variant it holds, so a variant goes wherever the union admitting it goes. A
// union owns no field, so this is the only edge leading out of one.
type VariantRule struct {
	variants parsed.Variants
}

// NewVariantRule constructs a VariantRule over a table of variants.
func NewVariantRule(variants parsed.Variants) VariantRule {
	return VariantRule{variants: variants}
}

// Edges implements [Rule].
func (r VariantRule) Edges() iter.Seq2[model.Reference, model.Reference] {
	return func(yield func(model.Reference, model.Reference) bool) {
		for key := range r.variants.All() {
			if !yield(key.Owner, key.Ref) {
				return
			}
		}
	}
}

// AliasRule is the [Rule] of standing for: an alias is another name and nothing
// besides, so the type it stands for goes wherever the alias goes. An alias
// standing for a primitive puts no edge in the graph, for the same reason a
// field typed as one does not.
type AliasRule struct {
	aliases flattened.Aliases
}

// NewAliasRule constructs an AliasRule over a table of aliases.
func NewAliasRule(aliases flattened.Aliases) AliasRule {
	return AliasRule{aliases: aliases}
}

// Edges implements [Rule].
func (r AliasRule) Edges() iter.Seq2[model.Reference, model.Reference] {
	return func(yield func(model.Reference, model.Reference) bool) {
		for ref, alias := range r.aliases.All() {
			named, ok := alias.Type.Atom().(typeform.Named)
			if !ok {
				continue
			}
			if !yield(ref, named.Ref()) {
				return
			}
		}
	}
}

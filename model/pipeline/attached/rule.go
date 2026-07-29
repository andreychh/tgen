// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package attached

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/flattened"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/typeform"
)

// Rule reports which definitions one relation of the specification adds to
// those already known to carry a file. A rule only reads the set it is given —
// the table that owns the set decides what to do with the answer — and may name
// a definition the set already holds.
type Rule interface {
	// Apply returns the definitions that carry a file because a definition in
	// files does.
	Apply(files Files) []model.Reference
}

// FieldRule is the [Rule] of ownership: a definition carries a file when a
// field it owns is typed as something that does. Dimensionality does not
// matter, since an array of carriers carries too.
type FieldRule struct {
	fields flattened.Fields
}

// NewFieldRule constructs a FieldRule over a table of fields.
func NewFieldRule(fields flattened.Fields) FieldRule {
	return FieldRule{fields: fields}
}

// Apply implements [Rule].
func (r FieldRule) Apply(files Files) []model.Reference {
	out := make([]model.Reference, 0)
	for key, field := range r.fields.All() {
		named, ok := field.Type.Atom().(typeform.Named)
		if !ok {
			continue
		}
		if _, carries := files.Lookup(named.Ref()); !carries {
			continue
		}
		out = append(out, key.Owner)
	}
	return out
}

// VariantRule is the [Rule] of alternation: a union carries a file when one of
// the definitions it is told apart into does. A union owns no field, so this is
// the only way a file reaches it.
type VariantRule struct {
	variants parsed.Variants
}

// NewVariantRule constructs a VariantRule over a table of variants.
func NewVariantRule(variants parsed.Variants) VariantRule {
	return VariantRule{variants: variants}
}

// Apply implements [Rule].
func (r VariantRule) Apply(files Files) []model.Reference {
	out := make([]model.Reference, 0)
	for key := range r.variants.All() {
		if _, carries := files.Lookup(key.Ref); !carries {
			continue
		}
		out = append(out, key.Owner)
	}
	return out
}

// AliasRule is the [Rule] of standing for: an alias carries a file when the
// type it stands for does.
type AliasRule struct {
	aliases flattened.Aliases
}

// NewAliasRule constructs an AliasRule over a table of aliases.
func NewAliasRule(aliases flattened.Aliases) AliasRule {
	return AliasRule{aliases: aliases}
}

// Apply implements [Rule].
func (r AliasRule) Apply(files Files) []model.Reference {
	out := make([]model.Reference, 0)
	for ref, alias := range r.aliases.All() {
		named, ok := alias.Type.Atom().(typeform.Named)
		if !ok {
			continue
		}
		if _, carries := files.Lookup(named.Ref()); !carries {
			continue
		}
		out = append(out, ref)
	}
	return out
}

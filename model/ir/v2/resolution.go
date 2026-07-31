// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"fmt"

	"github.com/andreychh/tgen/model/typebound"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/separated"
	"github.com/andreychh/tgen/model/typeform"
)

// Resolution is the join of a flat type with the tables naming its atom: it
// replaces the reference the atom addresses with the definition it names, and
// an alias with the type it stands for.
type Resolution struct {
	db  separated.Specification
	typ typeform.Type
}

// NewResolution constructs a Resolution of a flat type against the database
// naming its atom.
func NewResolution(db separated.Specification, typ typeform.Type) Resolution {
	return Resolution{db: db, typ: typ}
}

// Value returns the type with its atom resolved. It fails when the atom names
// no definition, names a method, or names an alias the aliases table does not
// hold.
func (r Resolution) Value() (typebound.Type, error) {
	switch atom := r.typ.Atom().(type) {
	case typeform.Primitive:
		return typebound.NewType(typebound.NewPrimitive(atom.Kind()), r.typ.Dimensionality()), nil
	case typeform.Named:
		resolved, err := r.named(atom.Ref())
		if err != nil {
			return typebound.Type{}, err
		}
		return typebound.NewType(resolved, r.typ.Dimensionality()), nil
	default:
		return typebound.Type{}, fmt.Errorf("resolving type atom: unexpected atom %T", atom)
	}
}

// named returns the atom naming the definition ref addresses.
func (r Resolution) named(ref model.Reference) (typebound.Atom, error) {
	definition, found := r.db.Definitions.Lookup(ref)
	if !found {
		return nil, fmt.Errorf("%s names no definition", ref)
	}
	switch definition.Kind {
	case model.DefinitionKindObject:
		return typebound.NewObject(definition.Name), nil
	case model.DefinitionKindUnion:
		return typebound.NewUnion(definition.Name), nil
	case model.DefinitionKindAlias:
		return r.alias(ref, definition.Name)
	case model.DefinitionKindMethod:
		return nil, fmt.Errorf("%s names a method, which is not a type", ref)
	default:
		return nil, fmt.Errorf("%s names a definition of unknown kind %q", ref, definition.Kind)
	}
}

// alias returns the atom naming the alias ref addresses, resolving the type it
// stands for.
func (r Resolution) alias(ref model.Reference, name model.Name) (typebound.Atom, error) {
	alias, found := r.db.Aliases.Lookup(ref)
	if !found {
		return nil, fmt.Errorf("alias %s stands for nothing", ref)
	}
	under, err := NewResolution(r.db, alias.Type).Value()
	if err != nil {
		return nil, fmt.Errorf("resolving what alias %s stands for: %w", ref, err)
	}
	return typebound.NewAlias(name, under), nil
}

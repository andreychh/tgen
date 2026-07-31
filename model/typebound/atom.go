// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package typebound

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/primitive"
)

// Atom represents what a type ultimately holds, with whatever it names already
// resolved. What a typeform atom leaves to a reference — the name of a
// definition and the shape a target must render it as — the variants carry
// themselves. The concrete variants are [Object], [Union], [Alias] and
// [Primitive].
//
//sumtype:decl
type Atom interface {
	isAtom()
}

// Object represents an atom naming an object: a definition with fields of its
// own.
type Object struct {
	name model.Name
}

// NewObject constructs an object atom from the name of the object it names.
func NewObject(name model.Name) Object {
	return Object{name: name}
}

// Name returns the name of the object the atom names.
func (o Object) Name() model.Name {
	return o.name
}

func (Object) isAtom() {}

// Union represents an atom naming a union: a definition standing for one of
// several variants.
type Union struct {
	name model.Name
}

// NewUnion constructs a union atom from the name of the union it names.
func NewUnion(name model.Name) Union {
	return Union{name: name}
}

// Name returns the name of the union the atom names.
func (u Union) Name() model.Name {
	return u.name
}

func (Union) isAtom() {}

// Alias represents an atom naming an alias: a name tgen gives a type the
// documentation leaves unnamed. It carries the type it stands for, since a
// target renders an alias by name but has to know what is underneath to spell
// its zero value.
type Alias struct {
	name  model.Name
	under Type
}

// NewAlias constructs an alias atom from the name of the alias and the type it
// stands for.
func NewAlias(name model.Name, under Type) Alias {
	return Alias{name: name, under: under}
}

// Name returns the name of the alias the atom names.
func (a Alias) Name() model.Name {
	return a.name
}

// Under returns the type the alias stands for.
func (a Alias) Under() Type {
	return a.under
}

func (Alias) isAtom() {}

// Primitive represents an atom naming a built-in type.
type Primitive struct {
	kind primitive.Kind
}

// NewPrimitive constructs a primitive atom of the given kind.
func NewPrimitive(kind primitive.Kind) Primitive {
	return Primitive{kind: kind}
}

// Kind returns the built-in type the atom names.
func (p Primitive) Kind() primitive.Kind {
	return p.kind
}

func (Primitive) isAtom() {}

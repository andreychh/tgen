// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golangv2

import (
	"fmt"
	"strings"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/primitive"
	"github.com/andreychh/tgen/model/typebound"
)

// Type represents a Go type expression rendered from a resolved type and the
// optionality of whatever carries it.
type Type struct {
	typ typebound.Type
	opt model.Optionality
}

// NewType creates a Type from a resolved type and its optionality.
func NewType(typ typebound.Type, opt model.Optionality) Type {
	return Type{typ: typ, opt: opt}
}

// NewRequiredType creates a Type from a resolved type no optionality applies
// to, such as the return of a method.
func NewRequiredType(typ typebound.Type) Type {
	return NewType(typ, false)
}

// Value returns the Go type expression: the name of the atom, enclosed in one
// slice per dimension, or behind a pointer when an optional single value needs
// one to tell an absent field from a zero one.
func (t Type) Value() string {
	if t.pointer() {
		return "*" + t.Name()
	}
	return strings.Repeat("[]", int(t.typ.Dimensionality())) + t.Name()
}

// Name returns the Go name of the atom the type holds, without the slices or
// the pointer enclosing it.
func (t Type) Name() string {
	switch atom := t.typ.Atom().(type) {
	case typebound.Primitive:
		return builtin(atom.Kind())
	case typebound.Object:
		return NewName(atom.Name()).Value()
	case typebound.Union:
		return NewName(atom.Name()).Value()
	case typebound.Alias:
		return NewName(atom.Name()).Value()
	default:
		panic(fmt.Sprintf("golangv2: unknown atom %T", atom))
	}
}

// Zero returns the Go expression the type's zero value is written as: the
// literal of a built-in, an empty composite of an object, what an alias stands
// for, and nil for everything Go already gives a nil.
func (t Type) Zero() string {
	if t.pointer() || t.Array() {
		return "nil"
	}
	switch atom := t.typ.Atom().(type) {
	case typebound.Primitive:
		return literal(atom.Kind())
	case typebound.Object:
		return NewName(atom.Name()).Value() + "{}"
	case typebound.Union:
		return "nil"
	case typebound.Alias:
		return NewRequiredType(atom.Under()).Zero()
	default:
		panic(fmt.Sprintf("golangv2: unknown atom %T", atom))
	}
}

// Union reports whether the atom the type holds names a union.
func (t Type) Union() bool {
	_, union := t.typ.Atom().(typebound.Union)
	return union
}

// Array reports whether the type encloses its atom in at least one slice.
func (t Type) Array() bool {
	return t.typ.Dimensionality() > 0
}

// pointer reports whether the type renders behind a pointer: only an optional
// single value does, since a slice and a union interface already carry nil.
func (t Type) pointer() bool {
	if !t.opt {
		return false
	}
	return !t.Array() && !t.Union()
}

// builtin returns the Go type a built-in of the documentation renders as.
func builtin(kind primitive.Kind) string {
	switch kind {
	case primitive.Integer:
		return "int64"
	case primitive.Float:
		return "float64"
	case primitive.String:
		return "string"
	case primitive.Boolean, primitive.True:
		return "bool"
	default:
		panic(fmt.Sprintf("golangv2: unknown primitive %q", kind))
	}
}

// literal returns the Go expression the zero value of a built-in is written as.
func literal(kind primitive.Kind) string {
	switch kind {
	case primitive.Integer, primitive.Float:
		return "0"
	case primitive.String:
		return `""`
	case primitive.Boolean, primitive.True:
		return "false"
	default:
		panic(fmt.Sprintf("golangv2: unknown primitive %q", kind))
	}
}

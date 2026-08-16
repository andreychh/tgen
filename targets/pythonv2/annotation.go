// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/primitive"
	"github.com/andreychh/tgen/model/typebound"
)

// Annotation represents a Python type annotation rendered from a resolved type
// and the optionality of whatever carries it.
type Annotation struct {
	typ typebound.Type
	opt model.Optionality
}

// NewAnnotation creates an Annotation from a resolved type and its optionality.
func NewAnnotation(typ typebound.Type, opt model.Optionality) Annotation {
	return Annotation{typ: typ, opt: opt}
}

// Value returns the annotation: the name of the atom, enclosed in one list per
// dimension, admitting None where whatever carries it is optional. An optional
// list and an optional single value are annotated alike, since None is a value
// the annotation admits rather than a way of holding one — which is what spares
// Python the pointer the Go target reaches for here.
func (a Annotation) Value() string {
	out := a.name()
	for range a.typ.Dimensionality() {
		out = "list[" + out + "]"
	}
	if a.opt {
		return out + " | None"
	}
	return out
}

// name returns the Python name of the atom the type holds, without the lists
// enclosing it.
func (a Annotation) name() string {
	switch atom := a.typ.Atom().(type) {
	case typebound.Primitive:
		return builtin(atom.Kind())
	case typebound.Object:
		return NewClassName(atom.Name()).Value()
	case typebound.Union:
		return NewClassName(atom.Name()).Value()
	case typebound.Alias:
		return NewClassName(atom.Name()).Value()
	}
	panic(fmt.Sprintf("pythonv2: unknown atom %T", a.typ.Atom()))
}

// builtin returns the Python type a built-in of the documentation renders as.
func builtin(kind primitive.Kind) string {
	switch kind {
	case primitive.Integer:
		return "int"
	case primitive.Float:
		return "float"
	case primitive.String:
		return "str"
	case primitive.Boolean, primitive.True:
		return "bool"
	default:
		panic(fmt.Sprintf("pythonv2: unknown primitive %q", kind))
	}
}

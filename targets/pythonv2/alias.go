// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	ir "github.com/andreychh/tgen/model/ir/v2"
)

// Alias represents the Python declaration of a name tgen gives a type the
// documentation leaves unnamed.
type Alias struct {
	inner ir.Alias
}

// NewAlias creates an Alias from the record of an alias.
func NewAlias(a ir.Alias) Alias {
	return Alias{inner: a}
}

// Doc returns the docstring of the declaration. An alias carries no link back
// to the documentation: tgen introduces it, so no section documents it.
func (a Alias) Doc() string {
	return NewClassDocstring(a.inner.Description).Value()
}

// Ref implements [Declaration].
func (a Alias) Ref() string {
	return string(a.inner.Ref)
}

// Template implements [Declaration].
func (a Alias) Template() string {
	return "alias"
}

// Name returns the Python class the alias declares.
func (a Alias) Name() string {
	return NewClassName(a.inner.Name).Value()
}

// Annotation returns the Python annotation the alias declares its name for.
func (a Alias) Annotation() string {
	return NewRequiredAnnotation(a.inner.Type).Value()
}

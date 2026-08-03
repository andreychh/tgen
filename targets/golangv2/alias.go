// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golangv2

import (
	ir "github.com/andreychh/tgen/model/ir/v2"
)

// Alias represents the Go declaration of a name tgen gives a type the
// documentation leaves unnamed.
type Alias struct {
	inner ir.Alias
}

// NewAlias creates an Alias from the record of an alias.
func NewAlias(a ir.Alias) Alias {
	return Alias{inner: a}
}

// Doc returns the doc comment of the declaration. An alias carries no link
// back to the documentation: tgen introduces it, so no section documents it.
func (a Alias) Doc() string {
	return NewTypeGodoc(a.inner.Description).Value()
}

// Ref implements [Declaration].
func (a Alias) Ref() string {
	return string(a.inner.Ref)
}

// Template implements [Declaration].
func (a Alias) Template() string {
	return "alias"
}

// Name returns the Go name the alias declares.
func (a Alias) Name() string {
	return NewName(a.inner.Name).Value()
}

// Type returns the Go type expression the alias declares its name for.
func (a Alias) Type() string {
	return NewRequiredType(a.inner.Type).Value()
}

// Direction returns which way the alias travels. The declaration itself is the
// same either way; what a block written by hand puts beside it is not.
func (a Alias) Direction() Direction {
	return NewDirection(a.inner.Direction)
}

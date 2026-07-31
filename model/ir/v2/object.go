// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typebound"
)

// Object is the record of a documented object, joined with the fields it owns.
// Files narrows Fields to those reaching a file, and is empty for an object
// holding none. Rewrites reports that a union reaching a file admits the object,
// which obliges it to rewrite itself into JSON even when it holds no file of its
// own: a union carries a file as soon as one variant does, and every target
// reaches the rest of them through the same rewrite.
type Object struct {
	Ref         model.Reference
	Name        model.Name
	Description prose.Passage
	Fields      []Field
	Files       []FileField
	Rewrites    bool
}

func (Object) isDefinition() {}

// Field is the record of a field an object owns or a parameter a method takes,
// with its type bound to the definition it names.
type Field struct {
	Key         model.Key
	Type        typebound.Type
	Optionality model.Optionality
	Description prose.Phrase
}

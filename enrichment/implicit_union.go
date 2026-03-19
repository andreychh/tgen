// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import (
	"iter"
	"slices"

	"github.com/andreychh/tgen/parsing"
)

type ImplicitUnion struct {
	name        string
	description string
	variants    []ImplicitVariant
}

func NewImplicitUnion(name, description string, variants []ImplicitVariant) ImplicitUnion {
	return ImplicitUnion{
		name:        name,
		description: description,
		variants:    variants,
	}
}

//nolint:ireturn // ObjectName is the intentional public contract of this method
func (u ImplicitUnion) Name() parsing.ObjectName {
	return staticObjectName{u.name}
}

//nolint:ireturn // DefinitionDescription is the intentional public contract of this method
func (u ImplicitUnion) Description() parsing.DefinitionDescription {
	return staticDefinitionDescription{u.description}
}

func (u ImplicitUnion) Variants() iter.Seq[ImplicitVariant] {
	return slices.Values(u.variants)
}

type staticObjectName struct {
	value string
}

func (n staticObjectName) Value() (string, error) {
	return n.value, nil
}

type staticDefinitionDescription struct {
	value string
}

func (s staticDefinitionDescription) Value() (string, error) {
	return s.value, nil
}

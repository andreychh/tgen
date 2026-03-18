// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
)

type ImplicitUnion struct {
	name        string
	description string
	variants    []string
}

func NewImplicitUnion(name, description string, variants []string) ImplicitUnion {
	return ImplicitUnion{
		name:        name,
		description: description,
		variants:    variants,
	}
}

func (u ImplicitUnion) Name() parsing.ObjectName {
	return staticObjectName{u.name}
}

func (u ImplicitUnion) Description() parsing.DefinitionDescription {
	return staticDefinitionDescription{u.description}
}

func (u ImplicitUnion) Variants() iter.Seq[parsing.Variant] {
	return func(yield func(parsing.Variant) bool) {
		for _, v := range u.variants {
			if !yield(staticVariant{name: v}) {
				break
			}
		}
	}
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

type staticVariant struct {
	name string
}

func (v staticVariant) Name() parsing.ObjectName {
	return staticObjectName{value: v.name}
}

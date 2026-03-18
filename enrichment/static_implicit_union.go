// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
)

// StaticImplicitUnion represents a parsing.ImplicitUnion defined by tgen
// from a hardcoded list of variant names.
type StaticImplicitUnion struct {
	name        string
	description string
	variants    []string
}

// NewStaticImplicitUnion creates a StaticImplicitUnion from a name, a
// description, and a list of variant type names.
func NewStaticImplicitUnion(name, description string, variants []string) StaticImplicitUnion {
	return StaticImplicitUnion{name: name, description: description, variants: variants}
}

func (u StaticImplicitUnion) Name() parsing.ObjectName {
	return staticObjectName{u.name}
}

func (u StaticImplicitUnion) Description() string {
	return u.description
}

func (u StaticImplicitUnion) Variants() iter.Seq[parsing.Variant] {
	return func(yield func(parsing.Variant) bool) {
		for _, v := range u.variants {
			if !yield(staticVariant{name: v}) {
				break
			}
		}
	}
}

type staticVariant struct {
	name string
}

func (v staticVariant) Name() parsing.ObjectName {
	return staticObjectName{value: v.name}
}

type staticObjectName struct {
	value string
}

func (n staticObjectName) Value() (string, error) {
	return n.value, nil
}

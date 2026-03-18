// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import "github.com/andreychh/tgen/parsing"

// FieldRule represents a transformation applied to a parsed field during
// enrichment.
type FieldRule interface {
	Apply(parsing.Field) parsing.Field
}

type typedField struct {
	inner parsing.Field
	tree  parsing.TypeTree
}

func (f typedField) Key() parsing.FieldKey {
	return f.inner.Key()
}

func (f typedField) IsOptional() parsing.Optionality {
	return f.inner.IsOptional()
}

func (f typedField) Description() parsing.FieldDescription {
	return f.inner.Description()
}

func (f typedField) Type() parsing.TypeTree {
	return f.tree
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import "github.com/andreychh/tgen/parsing/gq"

// GQVariantObject represents a variant of a discriminated union parsed from the
// specification.
type GQVariantObject struct {
	selection gq.Selection
}

// NewVariantObject constructs a GQVariantObject from an h4 selection.
func NewVariantObject(h4 gq.Selection) GQVariantObject {
	return GQVariantObject{selection: h4}
}

//nolint:ireturn // ObjectName is the intentional public contract of this method
func (v GQVariantObject) Name() ObjectName {
	return NewGQObjectName(v.selection)
}

//nolint:ireturn // VariantFields is the intentional public contract of this method
func (v GQVariantObject) Fields() VariantFields {
	return NewGQVariantFields(v.selection)
}

func (v GQVariantObject) Ref() DefinitionRef {
	return NewDefinitionRef(v.selection.Find("a.anchor"))
}

func (v GQVariantObject) Description() DefinitionDescription {
	return NewDefinitionDescription(v.selection)
}

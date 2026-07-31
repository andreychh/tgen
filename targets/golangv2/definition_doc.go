// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golangv2

import (
	"slices"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/targets"
)

// DefinitionDoc represents the doc comment of a documented definition: the
// prose of its section followed by a link back to that section.
type DefinitionDoc struct {
	ref     model.Reference
	passage prose.Passage
}

// NewDefinitionDoc creates a DefinitionDoc for the definition at ref from the
// prose of its section.
func NewDefinitionDoc(ref model.Reference, passage prose.Passage) DefinitionDoc {
	return DefinitionDoc{ref: ref, passage: passage}
}

// Value returns the doc comment, closing with the URL of the section.
func (d DefinitionDoc) Value() string {
	return NewTypeGodoc(
		prose.NewPassage(append(
			slices.Clone(d.passage.Blocks()),
			prose.NewParagraph(prose.NewText(
				"See "+targets.NewTelegramURL(d.ref).Value(),
				prose.StylePlain,
			)),
		)...),
	).Value()
}

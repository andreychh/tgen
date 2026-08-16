// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"slices"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/targets"
)

// DefinitionDoc represents the docstring of a definition: the prose describing
// it, closed by a link to the section of the documentation page it stands at,
// where it stands at one.
type DefinitionDoc struct {
	ref        model.Reference
	passage    prose.Passage
	introduced bool
}

// NewDefinitionDoc creates a DefinitionDoc for the definition at ref from the
// prose describing it and whether tgen introduced it.
func NewDefinitionDoc(ref model.Reference, passage prose.Passage, introduced bool) DefinitionDoc {
	return DefinitionDoc{ref: ref, passage: passage, introduced: introduced}
}

// Passage returns the prose describing the definition, closed by the URL of the
// section unless tgen introduced it, which the page never named and so gave no
// section to address. What indentation the prose is written at belongs to
// [Docstring], since it follows from where the declaration stands and not from
// what the documentation says.
func (d DefinitionDoc) Passage() prose.Passage {
	if d.introduced {
		return d.passage
	}
	return prose.NewPassage(append(
		slices.Clone(d.passage.Blocks()),
		prose.NewParagraph(prose.NewText(
			"See "+targets.NewTelegramURL(d.ref).Value(),
			prose.StylePlain,
		)),
	)...)
}

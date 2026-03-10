// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/parsing"
)

type DefinitionDoc struct {
	ref         parsing.DefinitionRef
	description parsing.DefinitionDescription
}

func NewDefinitionDoc(r parsing.DefinitionRef, d parsing.DefinitionDescription) DefinitionDoc {
	return DefinitionDoc{ref: r, description: d}
}

func (d DefinitionDoc) Value() (string, error) {
	ref, err := d.ref.Value()
	if err != nil {
		return "", fmt.Errorf("getting object ref: %w", err)
	}
	desc, err := d.description.Value()
	if err != nil {
		return "", fmt.Errorf("getting object description: %w", err)
	}
	return fmt.Sprintf("%s\n\nSee %s#%s", desc, "https://core.telegram.org/bots/api", ref), nil
}

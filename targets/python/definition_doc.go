// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/targets"
)

type DefinitionDoc struct {
	ref  model.Reference
	decs model.Description
}

func NewDefinitionDoc(r model.Reference, d model.Description) DefinitionDoc {
	return DefinitionDoc{ref: r, decs: d}
}

func (d DefinitionDoc) Value() (string, error) {
	desc, err := d.decs.Value()
	if err != nil {
		return "", fmt.Errorf("getting decs: %w", err)
	}
	return fmt.Sprintf("%s\n\nSee %s", desc, targets.NewTelegramURL(d.ref).Value()), nil
}

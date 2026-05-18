// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/targets"
)

type DefinitionDoc struct {
	ref  model.Reference
	desc model.Description
}

func NewDefinitionDoc(r model.Reference, d model.Description) DefinitionDoc {
	return DefinitionDoc{ref: r, desc: d}
}

func (d DefinitionDoc) Value() (string, error) {
	desc, err := d.desc.Value()
	if err != nil {
		return "", fmt.Errorf("getting description: %w", err)
	}
	return fmt.Sprintf("%s\n\nSee %s", desc, targets.NewTelegramURL(d.ref).Value()), nil
}

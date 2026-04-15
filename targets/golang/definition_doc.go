// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"fmt"

	"github.com/andreychh/tgen/model"
)

type DefinitionDoc struct {
	ref  model.Reference
	desc model.Description
}

func NewDefinitionDoc(r model.Reference, d model.Description) DefinitionDoc {
	return DefinitionDoc{ref: r, desc: d}
}

func (d DefinitionDoc) AsString() (string, error) {
	ref, err := d.ref.AsString()
	if err != nil {
		return "", fmt.Errorf("getting reference: %w", err)
	}
	desc, err := d.desc.AsString()
	if err != nil {
		return "", fmt.Errorf("getting description: %w", err)
	}
	return fmt.Sprintf("%s\n\nSee %s#%s", desc, specificationURL, ref), nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/targets"
)

type Doc struct {
	ref  model.Reference
	decs model.Description
}

func NewDoc(r model.Reference, d model.Description) Doc {
	return Doc{ref: r, decs: d}
}

func (d Doc) AsString() (string, error) {
	desc, err := d.decs.AsString()
	if err != nil {
		return "", fmt.Errorf("getting decs: %w", err)
	}
	url, err := targets.NewTelegramURL(d.ref).AsString()
	if err != nil {
		return "", fmt.Errorf("getting url: %w", err)
	}
	return fmt.Sprintf("%s\n\nSee %s", desc, url), nil
}

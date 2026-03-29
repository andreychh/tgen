// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"strings"

	"github.com/andreychh/tgen/parsing/gq"
)

type GQObjectFieldDescription struct {
	td gq.Selection
}

func NewGQObjectFieldDescription(td gq.Selection) GQObjectFieldDescription {
	return GQObjectFieldDescription{td: td}
}

func (d GQObjectFieldDescription) AsString() (string, error) {
	if d.td.IsEmpty() {
		return "", errors.New("description column not found")
	}
	return strings.TrimPrefix(d.td.Text(), "Optional. "), nil
}

func (d GQObjectFieldDescription) Links() ([]string, error) {
	if d.td.IsEmpty() {
		return nil, errors.New("description column not found")
	}
	var links []string
	for a := range d.td.Find("a").All() {
		if href, ok := a.Attr("href"); ok {
			links = append(links, href)
		}
	}
	return links, nil
}

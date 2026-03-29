// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"

	"github.com/andreychh/tgen/parsing/gq"
)

type GQMethodFieldDescription struct {
	td gq.Selection
}

func NewGQMethodFieldDescription(td gq.Selection) GQMethodFieldDescription {
	return GQMethodFieldDescription{td: td}
}

func (d GQMethodFieldDescription) AsString() (string, error) {
	if d.td.IsEmpty() {
		return "", errors.New("description column not found")
	}
	return d.td.Text(), nil
}

func (d GQMethodFieldDescription) Links() ([]string, error) {
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

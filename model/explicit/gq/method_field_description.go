// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"

	"github.com/andreychh/tgen/pkg/gq"
)

type MethodFieldDescription struct {
	td gq.Selection
}

func NewMethodFieldDescription(td gq.Selection) MethodFieldDescription {
	return MethodFieldDescription{td: td}
}

func (d MethodFieldDescription) AsString() (string, error) {
	if d.td.IsEmpty() {
		return "", errors.New("description column not found")
	}
	return d.td.Text(), nil
}

func (d MethodFieldDescription) Links() ([]string, error) {
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

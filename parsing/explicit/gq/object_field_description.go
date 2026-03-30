// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"strings"

	"github.com/andreychh/tgen/pkg/gq"
)

type ObjectFieldDescription struct {
	td gq.Selection
}

func NewObjectFieldDescription(td gq.Selection) ObjectFieldDescription {
	return ObjectFieldDescription{td: td}
}

func (d ObjectFieldDescription) AsString() (string, error) {
	if d.td.IsEmpty() {
		return "", errors.New("description column not found")
	}
	return strings.TrimPrefix(d.td.Text(), "Optional. "), nil
}

func (d ObjectFieldDescription) Links() ([]string, error) {
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

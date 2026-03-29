// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"strings"

	"github.com/andreychh/tgen/parsing/gq"
)

type GQDefinitionDescription struct {
	h4 gq.Selection
}

func NewGQDefinitionDescription(h4 gq.Selection) GQDefinitionDescription {
	return GQDefinitionDescription{h4: h4}
}

func (d GQDefinitionDescription) AsString() (string, error) {
	nodes := d.h4.
		Until("h3, h4, hr").
		Filter("p, blockquote")
	if nodes.IsEmpty() {
		return "", errors.New("definition description not found")
	}
	var parts []string
	for node := range nodes.All() {
		parts = append(parts, nodeText(node))
	}
	return strings.Join(parts, "\n\n"), nil
}

func nodeText(sel gq.Selection) string {
	var parts []string
	for child := range sel.Find("p, li").All() {
		text := strings.TrimSpace(child.Text())
		if text != "" {
			parts = append(parts, text)
		}
	}
	if len(parts) > 0 {
		return strings.Join(parts, "\n\n")
	}
	return strings.TrimSpace(sel.Text())
}

func (d GQDefinitionDescription) Links() ([]string, error) {
	nodes := d.h4.
		Until("h3, h4, hr").
		Filter("p, blockquote")
	if nodes.IsEmpty() {
		return nil, errors.New("definition description not found")
	}
	var links []string
	for a := range nodes.All() {
		if href, ok := a.Attr("href"); ok {
			links = append(links, href)
		}
	}
	return links, nil
}

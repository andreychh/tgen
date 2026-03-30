// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package gq implements the explicit interfaces by extracting data from the
// Telegram Bot API HTML page using goquery selections. It is the only package
// in the pipeline that reads raw HTML; assembled other packages work with the
// explicit interfaces.
package gq

import (
	"errors"
	"strings"

	"github.com/andreychh/tgen/pkg/gq"
)

type DefinitionDescription struct {
	h4 gq.Selection
}

func NewDefinitionDescription(h4 gq.Selection) DefinitionDescription {
	return DefinitionDescription{h4: h4}
}

func (d DefinitionDescription) AsString() (string, error) {
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

func (d DefinitionDescription) Links() ([]string, error) {
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

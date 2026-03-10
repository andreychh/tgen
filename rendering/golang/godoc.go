// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"strings"

	"github.com/mitchellh/go-wordwrap"
)

type RawDesc interface {
	Value() (string, error)
}

type Godoc struct {
	inner RawDesc
}

func NewGodoc(d RawDesc) Godoc {
	return Godoc{inner: d}
}

func (g Godoc) Value() (string, error) {
	text, err := g.inner.Value()
	if err != nil {
		return "", err
	}
	paragraphs := strings.Split(text, "\n\n")
	var lines []string
	for i, p := range paragraphs {
		wrapped := wordwrap.WrapString(p, 77)
		for _, line := range strings.Split(wrapped, "\n") {
			lines = append(lines, "// "+line)
		}
		if i < len(paragraphs)-1 {
			lines = append(lines, "//")
		}
	}
	return strings.Join(lines, "\n"), nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"strings"

	"github.com/mitchellh/go-wordwrap"
)

//nolint:iface // intentionally distinct from Type: Type is a Go type expression, not a doc string
type Stringable interface {
	AsString() (string, error)
}

type GoDoc struct {
	inner Stringable
}

func NewGoDoc(s Stringable) GoDoc {
	return GoDoc{inner: s}
}

func (d GoDoc) AsString() (string, error) {
	text, err := d.inner.AsString()
	if err != nil {
		return "", err
	}
	paragraphs := strings.Split(text, "\n\n")
	var lines []string
	for i, p := range paragraphs {
		wrapped := wordwrap.WrapString(p, 77)
		for line := range strings.SplitSeq(wrapped, "\n") {
			lines = append(lines, "// "+line)
		}
		if i < len(paragraphs)-1 {
			lines = append(lines, "//")
		}
	}
	return strings.Join(lines, "\n"), nil
}

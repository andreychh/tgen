// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"strings"

	"github.com/andreychh/tgen/model"
	"github.com/mitchellh/go-wordwrap"
)

type Stringable interface {
	AsString() (string, error)
}

type DocString struct {
	description Stringable
	indent      int
}

func NewDocString(d Stringable, indent int) DocString {
	return DocString{description: d, indent: indent}
}

func NewFieldDocString(d model.Description) DocString {
	return NewDocString(d, 4)
}

func NewClassDocString(r model.Reference, d model.Description) DocString {
	return NewDocString(NewDoc(r, d), 4)
}

const (
	docLimit  = 72
	docQuotes = `"""`
)

func (d DocString) AsString() (string, error) {
	desc, err := d.description.AsString()
	if err != nil {
		return "", err
	}
	if len(desc) <= docLimit {
		return docQuotes + desc + docQuotes, nil
	}
	return d.multiline(desc), nil
}

func (d DocString) multiline(desc string) string {
	pad := strings.Repeat(" ", d.indent)
	paragraphs := strings.Split(desc, "\n\n")
	var out strings.Builder
	out.WriteString(docQuotes + "\n")
	for i, para := range paragraphs {
		wrapped := wordwrap.WrapString(para, docLimit)
		for line := range strings.SplitSeq(wrapped, "\n") {
			out.WriteString(pad + line + "\n")
		}
		if i < len(paragraphs)-1 {
			out.WriteString("\n")
		}
	}
	out.WriteString(pad + docQuotes)
	return out.String()
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"strings"

	"github.com/mitchellh/go-wordwrap"
)

type DocString struct {
	content string
	indent  int
}

func NewDocString(d string, indent int) DocString {
	return DocString{content: d, indent: indent}
}

func NewFieldDocString(d string) DocString {
	return NewDocString(d, 4)
}

func NewClassDocString(d string) DocString {
	return NewDocString(d, 4)
}

const (
	docLimit  = 72
	docQuotes = `"""`
)

func (d DocString) Value() string {
	if len(d.content) <= docLimit {
		return docQuotes + d.content + docQuotes
	}
	return d.multiline(d.content)
}

func (d DocString) multiline(desc string) string {
	pad := strings.Repeat(" ", d.indent)
	paragraphs := strings.Split(desc, "\n\n")
	var out strings.Builder
	out.WriteString(docQuotes + "\n")
	for i, par := range paragraphs {
		wrapped := wordwrap.WrapString(par, docLimit)
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

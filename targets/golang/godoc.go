// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"strings"

	"github.com/mitchellh/go-wordwrap"
)

type GoDoc struct {
	content string
	indent  int
}

func NewGoDoc(content string, indent int) GoDoc {
	return GoDoc{content: content, indent: indent}
}

func NewTypeGodoc(content string) GoDoc {
	return NewGoDoc(content, 0)
}

func NewFieldGodoc(content string) GoDoc {
	return NewGoDoc(content, 1)
}

func (d GoDoc) Value() string {
	pad := strings.Repeat("\t", d.indent)
	paragraphs := strings.Split(d.content, "\n\n")
	var lines []string
	for i, p := range paragraphs {
		wrapped := wordwrap.WrapString(p, 77)
		for line := range strings.SplitSeq(wrapped, "\n") {
			lines = append(lines, pad+"// "+line)
		}
		if i < len(paragraphs)-1 {
			lines = append(lines, pad+"//")
		}
	}
	return strings.Join(lines, "\n")
}

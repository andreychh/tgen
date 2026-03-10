// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"strings"

	"github.com/mitchellh/go-wordwrap"
)

type Doc struct {
	inner RawValue
}

func NewDoc(d RawValue) Doc {
	return Doc{inner: d}
}

func (d Doc) Value() (string, error) {
	text, err := d.inner.Value()
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

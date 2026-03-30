// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package literals

// Description represents a parsing.Description wrapping known text and links.
type Description struct {
	text  string
	links []string
}

// NewDescription constructs a Description from text and links.
func NewDescription(text string, links []string) Description {
	return Description{text: text, links: links}
}

func (d Description) AsString() (string, error) {
	return d.text, nil
}

func (d Description) Links() ([]string, error) {
	return d.links, nil
}

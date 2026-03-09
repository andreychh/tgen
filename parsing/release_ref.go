// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andreychh/tgen/parsing/gq"
)

var releaseRefRegex = regexp.MustCompile(`^#[a-z]+-\d+-\d+$`)

type ReleaseRef struct {
	selection gq.Selection
}

func NewReleaseRef(a gq.Selection) ReleaseRef {
	return ReleaseRef{selection: a}
}

func (r ReleaseRef) Value() (string, error) {
	val, _ := r.selection.Attr("href")
	if !releaseRefRegex.MatchString(val) {
		return "", fmt.Errorf("invalid release ref %q", val)
	}
	return strings.TrimPrefix(val, "#"), nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/andreychh/tgen/parsing/gq"
)

var releaseRefRegex = regexp.MustCompile(`^#[a-z]+-\d+-\d+$`)

type GQReleaseReference struct {
	a gq.Selection
}

func NewGQReleaseReference(a gq.Selection) GQReleaseReference {
	return GQReleaseReference{a: a}
}

func (r GQReleaseReference) AsString() (string, error) {
	if r.a.IsEmpty() {
		return "", errors.New("release reference not found")
	}
	val, _ := r.a.Attr("href")
	if !releaseRefRegex.MatchString(val) {
		return "", fmt.Errorf("release reference %q contains invalid characters", val)
	}
	return strings.TrimPrefix(val, "#"), nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/andreychh/tgen/pkg/gq"
)

var releaseRefRegex = regexp.MustCompile(`^#[a-z]+-\d+-\d+$`)

type ReleaseReference struct {
	a gq.Selection
}

func NewReleaseReference(a gq.Selection) ReleaseReference {
	return ReleaseReference{a: a}
}

func (r ReleaseReference) AsString() (string, error) {
	if r.a.IsEmpty() {
		return "", errors.New("release reference not found")
	}
	val, _ := r.a.Attr("href")
	if !releaseRefRegex.MatchString(val) {
		return "", fmt.Errorf("release reference %q contains invalid characters", val)
	}
	return strings.TrimPrefix(val, "#"), nil
}

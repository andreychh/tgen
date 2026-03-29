// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var releaseVersionRegex = regexp.MustCompile(`^Bot API (\d+\.\d+)$`)

type GQReleaseVersion struct {
	strong gq.Selection
}

func NewGQReleaseVersion(strong gq.Selection) GQReleaseVersion {
	return GQReleaseVersion{strong: strong}
}

func (v GQReleaseVersion) AsString() (string, error) {
	if v.strong.IsEmpty() {
		return "", errors.New("release version not found")
	}
	val := v.strong.Text()
	matches := releaseVersionRegex.FindStringSubmatch(val)
	if len(matches) != 2 {
		return "", fmt.Errorf("release version %q contains invalid characters", val)
	}
	return matches[1], nil
}

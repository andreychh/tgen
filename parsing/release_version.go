// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/parsing/gq"
)

var releaseVersionRegex = regexp.MustCompile(`^Bot API (\d+\.\d+)$`)

type ReleaseVersion struct {
	selection gq.Selection
}

func NewReleaseVersion(strong gq.Selection) ReleaseVersion {
	return ReleaseVersion{selection: strong}
}

func (v ReleaseVersion) Value() (string, error) {
	val := v.selection.Text()
	matches := releaseVersionRegex.FindStringSubmatch(val)
	if len(matches) != 2 {
		return "", fmt.Errorf("invalid release version: %q", val)
	}
	return matches[1], nil
}

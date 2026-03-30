// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/andreychh/tgen/pkg/gq"
)

var releaseVersionRegex = regexp.MustCompile(`^Bot API (\d+\.\d+)$`)

type ReleaseVersion struct {
	strong gq.Selection
}

func NewReleaseVersion(strong gq.Selection) ReleaseVersion {
	return ReleaseVersion{strong: strong}
}

func (v ReleaseVersion) AsString() (string, error) {
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

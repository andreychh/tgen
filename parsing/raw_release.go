// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"time"

	"github.com/andreychh/tgen/parsing/dom"
)

// RawRelease represents a parsed Telegram API release definition.
//
// It expects the selection to point to a release header (h4), immediately
// followed by version paragraphs.
type RawRelease struct {
	selection dom.Selection
}

// NewRawRelease creates a RawRelease starting from the given semantic DOM node.
func NewRawRelease(h4 dom.Selection) RawRelease {
	return RawRelease{selection: h4}
}

// ID returns the unique reference identifier of the release (e.g.,
// "#february-9-2026").
func (r RawRelease) ID() (string, error) {
	val, exists := r.selection.Find("a.anchor").Attr("href")
	if !exists {
		return "", errors.New("attribute href not found")
	}
	if !releaseIDRegex.MatchString(val) {
		return "", fmt.Errorf("id %q does not match pattern %s", val, releaseIDRegex)
	}
	return val, nil
}

// Version returns the API version identifier (e.g., "v9.4").
func (r RawRelease) Version() (string, error) {
	sel := r.selection.NextAllFiltered("p").Find("strong").First()
	if sel.IsEmpty() {
		return "", errors.New("API version element not found")
	}
	version := sel.Text()
	matches := versionRegex.FindStringSubmatch(version)
	if len(matches) != 2 {
		return "", fmt.Errorf("text %q does not match pattern %q", version, versionRegex)
	}
	return "v" + matches[1], nil
}

// Date parses and returns the release date.
func (r RawRelease) Date() (time.Time, error) {
	date, exists := r.selection.Find("a.anchor").Attr("href")
	if !exists {
		return time.Time{}, errors.New("attribute href not found")
	}
	parsed, err := time.Parse("#January-2-2006", date)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing release date: %w", err)
	}
	return parsed, nil
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsed

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/andreychh/tgen/model"
)

// refPattern matches a reference once its leading '#' is stripped.
var refPattern = regexp.MustCompile(`^[a-z0-9]+$`)

// Reference is the anchor of a section's <h4> header, addressing the section on
// the page.
type Reference struct {
	h4 *goquery.Selection
}

// NewReference constructs a Reference over a section's <h4> header.
func NewReference(h4 *goquery.Selection) Reference {
	return Reference{h4: h4}
}

// Value returns the reference with its leading '#' stripped. It fails when the
// anchor is absent or the reference is malformed.
func (r Reference) Value() (model.Reference, error) {
	href, found := r.h4.Find("a.anchor").Attr("href")
	if !found {
		return "", errors.New("anchor href not found")
	}
	ref := strings.TrimPrefix(href, "#")
	if !refPattern.MatchString(ref) {
		return "", fmt.Errorf("reference %q is malformed", ref)
	}
	return model.Reference(ref), nil
}

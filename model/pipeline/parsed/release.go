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

// releaseRefPattern matches the dated anchor of the latest release entry.
var releaseRefPattern = regexp.MustCompile(`^#[a-z]+-\d+-\d+$`)

// releaseVersionPattern captures the Bot API version from its strong element.
var releaseVersionPattern = regexp.MustCompile(`^Bot API (\d+\.\d+)$`)

// Release is the decoded record of a Bot API release: its reference and
// version. Unlike the rest of the decoded data it is a single value, not a
// table. The release date stays encoded in the reference for a later pass to
// lift.
type Release struct {
	Ref     model.Reference
	Version model.ReleaseVersion
}

// Changelog is the recent-changes section of a documentation page. Its entries
// are the spec's releases, latest first.
type Changelog struct {
	doc *goquery.Document
}

// NewChangelog constructs a Changelog over a parsed documentation page.
func NewChangelog(doc *goquery.Document) Changelog {
	return Changelog{doc: doc}
}

// Latest returns the most recent release: its reference and version. It fails
// when either is absent or malformed.
func (c Changelog) Latest() (Release, error) {
	head := c.doc.Find("div#dev_page_content h4").First()
	href, found := head.Find("a.anchor").Attr("href")
	if !found {
		return Release{}, errors.New("release reference not found")
	}
	if !releaseRefPattern.MatchString(href) {
		return Release{}, fmt.Errorf("release reference %q is malformed", href)
	}
	strong := head.Next().Find("strong").First()
	version := releaseVersionPattern.FindStringSubmatch(strong.Text())
	if version == nil {
		return Release{}, fmt.Errorf("release version %q is malformed", strong.Text())
	}
	return Release{
		Ref:     model.Reference(strings.TrimPrefix(href, "#")),
		Version: model.ReleaseVersion(version[1]),
	}, nil
}

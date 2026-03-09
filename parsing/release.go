// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"github.com/andreychh/tgen/parsing/gq"
)

type Release struct {
	selection gq.Selection
}

func NewRelease(h4 gq.Selection) Release {
	return Release{selection: h4}
}

func (r Release) Ref() ReleaseRef {
	return NewReleaseRef(r.selection.Find("a.anchor"))
}

func (r Release) Version() ReleaseVersion {
	return NewReleaseVersion(
		r.selection.
			Until("h3, h4, hr").
			Filter("p").
			Find("strong").
			At(0),
	)
}

func (r Release) Date() ReleaseDate {
	return NewReleaseDate(r.selection.Find("a.anchor"))
}

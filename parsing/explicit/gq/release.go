// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"github.com/andreychh/tgen/parsing/explicit"
	"github.com/andreychh/tgen/pkg/gq"
)

type Release struct {
	h4 gq.Selection
}

func NewRelease(h4 gq.Selection) Release {
	return Release{h4: h4}
}

func (r Release) Reference() explicit.Reference {
	return NewReleaseReference(r.h4.Find("a.anchor"))
}

func (r Release) Version() explicit.ReleaseVersion {
	return NewReleaseVersion(r.h4.
		Until("h3, h4, hr").
		Filter("p").
		Find("strong").
		At(0),
	)
}

func (r Release) Date() explicit.ReleaseDate {
	return NewReleaseDate(r.h4.Find("a.anchor"))
}

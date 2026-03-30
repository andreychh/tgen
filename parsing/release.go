// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import "github.com/andreychh/tgen/parsing/gq"

type GQRelease struct {
	h4 gq.Selection
}

func NewGQRelease(h4 gq.Selection) GQRelease {
	return GQRelease{h4: h4}
}

func (r GQRelease) Reference() Reference {
	return NewGQReleaseReference(r.h4.Find("a.anchor"))
}

func (r GQRelease) Version() ReleaseVersion {
	return NewGQReleaseVersion(r.h4.
		Until("h3, h4, hr").
		Filter("p").
		Find("strong").
		At(0),
	)
}

func (r GQRelease) Date() ReleaseDate {
	return NewGQReleaseDate(r.h4.Find("a.anchor"))
}

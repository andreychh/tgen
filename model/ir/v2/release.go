// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"github.com/andreychh/tgen/model"
)

// Release is the record of the Bot API release the specification was read
// from: the reference of the changelog entry announcing it, and the version
// that entry names. Unlike everything else a target reads here, it is a single
// value rather than a table, since a page announces one latest release.
type Release struct {
	Ref     model.Reference
	Version model.ReleaseVersion
}

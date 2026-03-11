// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package meta

// The following variables are intended to be populated at build time via
// -ldflags.
//
//nolint:gochecknoglobals // These variables are set at build time via ldflags
var (
	version   = "unknown"
	commit    = "unknown"
	date      = "unknown"
	treeState = "unknown"
	builder   = "unknown"
)

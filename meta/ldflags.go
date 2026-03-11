// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package meta

const unknownValue = "unknown"

// The following variables are intended to be populated at build time via
// -ldflags.
//
//nolint:gochecknoglobals // These variables are set at build time via ldflags
var (
	version   = unknownValue
	commit    = unknownValue
	date      = unknownValue
	treeState = unknownValue
	builder   = unknownValue
)

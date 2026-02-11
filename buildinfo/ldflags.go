// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package buildinfo provides access to binary metadata set during the build
// process.
package buildinfo

// The following variables are intended to be populated at build time via
// -ldflags.
//
//nolint:gochecknoglobals // These variables are set at build time via ldflags
var (
	version   = "dev"
	commit    = "unknown"
	branch    = "unknown"
	buildTime = "unknown"
	treeState = "unknown"
	builtBy   = "manual"
)

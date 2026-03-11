// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package meta

// Well-known keys for use with [Source.Get].
const (
	// KeyVersion is the semantic version of the release (e.g. "v1.2.3").
	KeyVersion = "version"

	// KeyBuilder is the name of the tool or scenario that produced the binary (e.g.
	// "goreleaser", "gotoolchain").
	KeyBuilder = "builder"

	// KeyVCSRevision is the full VCS commit hash.
	KeyVCSRevision = "vcs.revision"

	// KeyVCSTime is the RFC3339 commit timestamp.
	KeyVCSTime = "vcs.time"

	// KeyVCSTreeState indicates whether the working tree was dirty at build time.
	// The value is "true" if dirty, "false" if clean.
	KeyVCSTreeState = "vcs.modified"

	// KeyGoVersion is the version of the Go toolchain (e.g. "go1.26.1").
	KeyGoVersion = "go.version"

	// KeyPlatform is the target OS and architecture (e.g. "linux/amd64").
	KeyPlatform = "platform"
)

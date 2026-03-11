// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package cli

import (
	"fmt"
	"time"

	"github.com/andreychh/tgen/meta"
)

// VersionMessage formats binary metadata for the --version flag output.
type VersionMessage struct {
	meta meta.Meta
}

// NewVersionMessage creates a VersionMessage.
func NewVersionMessage(m meta.Meta) VersionMessage {
	return VersionMessage{meta: m}
}

// String returns the formatted version string.
func (v VersionMessage) String() string {
	vcs := v.meta.VCS()
	release := v.meta.Release()
	build := v.meta.Build()
	return fmt.Sprintf(`tgen %s
https://github.com/andreychh/tgen

commit:      %s
date:        %s
tree state:  %s
built by:    %s
go version:  %s
platform:    %s
`,
		release.Version(),
		vcs.Revision().Full(),
		vcs.Date().Format(time.RFC3339),
		vcs.TreeState(),
		release.Builder(),
		build.GoVersion(),
		build.Platform(),
	)
}

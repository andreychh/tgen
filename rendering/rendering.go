// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package rendering

import (
	"io"
)

type View interface {
	Render(w io.Writer) error
}

type Artifacts map[string]View

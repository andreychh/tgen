// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"github.com/andreychh/tgen/model"
)

// FileField is a field that reaches a file, joined with the way it reaches one.
// It is a field of its owner narrowed to one concern, so it carries everything
// [Field] carries and never stands on its own.
type FileField struct {
	Field

	Kind model.FileKind
}

// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package overlays applies editorial type corrections to parsed Telegram Bot
// API fields and methods.
package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/pkg/iters"
)

// Overlay represents a conditional field transformation applied during
// parsing.
type Overlay interface {
	Apply(f spec.Field) spec.Field
}

func NewPrioritizedFields(seq iter.Seq[spec.Field]) iter.Seq[spec.Field] {
	return iters.NewMergedSeq(
		iters.NewFilteredSeq(seq, func(f spec.Field) bool {
			opt, err := f.Optionality()
			if err != nil {
				return true
			}
			return !bool(opt)
		}),
		iters.NewFilteredSeq(seq, func(f spec.Field) bool {
			opt, err := f.Optionality()
			if err != nil {
				return false
			}
			return bool(opt)
		}),
	)
}

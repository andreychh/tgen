// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package overlays applies editorial type corrections to parsed Telegram Bot
// API fields and methods.
package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// Overlay represents a conditional field transformation applied during
// parsing.
type Overlay interface {
	Apply(f explicit.Field) explicit.Field
}

func NewPrioritizedFields(seq iter.Seq[explicit.Field]) iter.Seq[explicit.Field] {
	return iters.NewMergedSeq(
		iters.NewFilteredSeq(seq, func(f explicit.Field) bool {
			opt, err := f.Optionality().AsBool()
			if err != nil {
				return true
			}
			return !opt
		}),
		iters.NewFilteredSeq(seq, func(f explicit.Field) bool {
			opt, err := f.Optionality().AsBool()
			if err != nil {
				return false
			}
			return opt
		}),
	)
}

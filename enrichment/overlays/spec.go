// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
)

// Spec represents the Telegram Bot API specification with field and method type overlays applied.
type Spec struct {
	inner   parsing.Specification
	overlay Overlay
}

// NewSpec constructs a Spec from a parsed specification.
func NewSpec(s parsing.Specification) Spec {
	return Spec{
		inner: s,
		overlay: NewSequential(
			ChatID{},
			ReplyMarkup{},
			InputMediaGroup{},
			InputFile{},
		),
	}
}

func (s Spec) Objects() iter.Seq[Object] {
	return func(yield func(Object) bool) {
		for o := range s.inner.Objects() {
			if !yield(NewObject(o, s.overlay)) {
				break
			}
		}
	}
}

func (s Spec) Methods() iter.Seq[Method] {
	return func(yield func(Method) bool) {
		for m := range s.inner.Methods() {
			if !yield(NewMethod(m, s.overlay)) {
				break
			}
		}
	}
}

func (s Spec) Unions() parsing.Unions   { return s.inner.Unions() }
func (s Spec) Release() parsing.Release { return s.inner.Release() }
